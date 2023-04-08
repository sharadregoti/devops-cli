import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { DispatchType, StoreType } from "./store";

export const useAppDispatch: () => DispatchType = useDispatch;
export const useAppSelector: TypedUseSelectorHook<StoreType> = useSelector;