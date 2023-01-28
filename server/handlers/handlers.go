package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	core "github.com/sharadregoti/devops"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
)

// HandleInfo gives info
// @Summary      HandleInfo endpoint
// @ID           HandleInfo
// @Description  HandleInfo endpoint
// @Accept       json
// @Produce      json
// @Success      200 {object} model.Info true
// @Failure      400 {object} model.ErrorResponse true
// @Failure      404 {object} model.ErrorResponse true
// @Failure      500 {object} model.ErrorResponse true
// @Router       /v1/info [get]
func HandleInfo(sm *core.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(5)*time.Second)
		defer cancel()

		ID := fmt.Sprintf("%d", sm.SessionCount()+1)
		if err := sm.AddClient(ID); err != nil {
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		pCtx, err := sm.GetClient(ID)
		if err != nil {
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		_ = utils.SendResponse(ctx, w, http.StatusOK, pCtx.GetInfo(ID))
	}
}

// HandleEvent handles events
// @Summary      HandleEvent endpoint
// @ID           HandleEvent
// @Description  HandleEvent endpoint
// @Accept       json
// @Produce      json
// @Success      200 {object} model.EventResponse true
// @Failure      400 {object} model.ErrorResponse true
// @Failure      404 {object} model.ErrorResponse true
// @Failure      500 {object} model.ErrorResponse true
// @Router       /v1/events/{id} [post]
func HandleEvent(sm *core.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(5)*time.Second)
		defer cancel()

		req := new(model.FrontendEvent)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		defer utils.CloseTheCloser(r.Body)

		params := mux.Vars(r)
		ID := params["id"]
		pCtx, err := sm.GetClient(ID)
		if err != nil {
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		if strings.ToLower(req.IsolatorName) == "na" {
			req.IsolatorName = ""
		}
		e := model.Event{
			Type:         string(req.ActionName),
			ResourceName: req.ResourceName,
			ResourceType: req.ResourceType,
			IsolatorName: req.IsolatorName,
			PluginName:   req.PluginName,
			Args:         req.Args,
			// TODO: Remove this
			// RowIndex: "",
			// SpecificActionName: "",
		}

		dd, _ := json.MarshalIndent(e, " ", " ")

		logger.LogDebug("Received event \n(%v)", string(dd))

		switch req.EventType {
		case model.NormalAction:
			switch model.NormalEvent(req.ActionName) {

			case model.CreateResource:

			case model.EditResource:

			case model.ReadResource:
				result, err := pCtx.Read(e)
				if err != nil {
					_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
					return
				}

				_ = utils.SendResponse(r.Context(), w, http.StatusOK, &model.EventResponse{Result: result})
				return

			case model.DeleteResource:
				if pCtx.Delete(e); err != nil {
					_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
					return
				}

				_ = utils.SendResponse(r.Context(), w, http.StatusOK, &model.EventResponse{})
				return

			case model.IsolatorChanged:
				if pCtx.ReadSync(e); err != nil {
					_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
					return
				}

				_ = utils.SendResponse(r.Context(), w, http.StatusOK, &model.EventResponse{})
				return

			default:
				_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, fmt.Errorf("unknown normal event  %v provided", req.EventType))
				return
			}
		case model.InternalAction:
			switch model.InternalEvent(req.ActionName) {

			case model.RefreshResource:
				if pCtx.ReadSync(e); err != nil {
					_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
					return
				}

				_ = utils.SendResponse(r.Context(), w, http.StatusOK, &model.EventResponse{})
				return

			case model.ResourceTypeChanged:
				if pCtx.ReadSync(e); err != nil {
					_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
					return
				}

				_ = utils.SendResponse(r.Context(), w, http.StatusOK, &model.EventResponse{})
				return

			default:
				_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, fmt.Errorf("unknown internal event  %v provided", req.EventType))
				return
			}
		case model.SpecificAction:
			actions, err := pCtx.GetSpecficActionList(e)
			if err != nil {
				_ = utils.SendResponse(r.Context(), w, http.StatusOK, &model.EventResponse{})
				return
			}

			for _, a := range actions {
				if a.Name == req.ActionName {
					logger.LogDebug("Specific action match found")

					res, err := pCtx.SpecificAction(a, e)
					if err != nil {
						_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
						return
					}

					_ = utils.SendResponse(r.Context(), w, http.StatusOK, res)
					return
				}
			}

			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, fmt.Errorf("action not found"))
			return

		default:
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, fmt.Errorf("invalid event type %v provided", req.EventType))
			return
		}

		_ = utils.SendResponse(ctx, w, http.StatusOK, map[string]string{})
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     utils.CreateCorsObject().OriginAllowed,
}

func HandleWebsocket(sm *core.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		defer conn.Close()

		conn.SetCloseHandler(func(code int, text string) error {
			logger.LogInfo("Websocket connection closed with code (%v) text (%v)", code, text)
			return err
		})

		params := mux.Vars(r)
		ID := params["id"]

		pCtx, err := sm.GetClient(ID)
		if err != nil {
			data, _ := json.Marshal(model.ErrorResponse{Message: err.Error()})
			conn.WriteMessage(websocket.CloseMessage, data)
			return
		}

		go func() {
			_, _, err := conn.ReadMessage()
			if err != nil {
				sm.DeleteClient(ID)
			}
		}()

		for v := range pCtx.GetDataPipe() {
			if err = conn.WriteJSON(v); err != nil {
				data, _ := json.Marshal(model.ErrorResponse{Message: err.Error()})
				conn.WriteMessage(websocket.CloseMessage, data)
				return
			}
		}
	}
}

type socketReadWriter struct {
	Socket *websocket.Conn
}

func (srw socketReadWriter) Write(p []byte) (n int, err error) {
	err = srw.Socket.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

func (srw socketReadWriter) Read(p []byte) (n int, err error) {
	_, b, err := srw.Socket.ReadMessage()
	for i, d := range b {
		p[i] = d
	}
	return len(b), err
}

func HandleActionWebsocket(sm *core.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_ = utils.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		defer conn.Close()

		params := mux.Vars(r)
		clientID := params["clientId"]
		ID := params["id"]

		pCtx, err := sm.GetClient(clientID)
		if err != nil {
			data, _ := json.Marshal(model.ErrorResponse{Message: err.Error()})
			conn.WriteMessage(websocket.CloseMessage, data)
			return
		}

		// Invoke specific event
		if err := pCtx.PerformSavedAction(ID, socketReadWriter{Socket: conn}); err != nil {
			data, _ := json.Marshal(model.ErrorResponse{Message: fmt.Errorf("failed to perform specific action: %v", err).Error()})
			conn.WriteMessage(websocket.CloseMessage, data)
			return
		}

		// // Setup reader & writer chans
		// readerChan := make(chan string, 1)
		// writerChan := make(chan string, 1)
		// defer close(readerChan)
		// defer close(writerChan)

		// // TODO: Add closer for this routines
		// pCtx.ReadFromSTDOUT(readerChan)
		// pCtx.WriteIntoSTDIN(writerChan)

		// blocker := make(chan struct{}, 1)
		// go func() {
		// 	for {
		// 		// Write into stdin
		// 		_, data, err := conn.ReadMessage()
		// 		if err != nil {
		// 			logger.LogError("Failed to read message: %v", err)
		// 			conn.Close()
		// 			blocker <- struct{}{}
		// 			return
		// 		}

		// 		logger.LogDebug("Recived data")
		// 		writerChan <- string(data)
		// 		logger.LogDebug("Sent receive data to chan")
		// 	}
		// }()

		// go func() {
		// 	// Read from stdout
		// 	for v := range readerChan {
		// 		if err = conn.WriteMessage(websocket.TextMessage, []byte(v)); err != nil {
		// 			// TODO: Check how error will be delivered
		// 			logger.LogError("Failed to write message: %v", err)
		// 			conn.Close()
		// 			return
		// 		}
		// 	}
		// }()

	}
}
