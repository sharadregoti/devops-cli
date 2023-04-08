import React, { useRef, useEffect, useLayoutEffect } from 'react';
import { Terminal } from 'xterm';
import { AttachAddon } from 'xterm-addon-attach';
import { FitAddon } from 'xterm-addon-fit';

type XTermComponentPropsTypes = {
    socketUrl: string,
    isDrawerOpen: boolean
}

const XTermComponent: React.FC<XTermComponentPropsTypes> = ({ isDrawerOpen, socketUrl }) => {
    const containerRef = useRef(null);
    const terminalRef = useRef(null);

    const fitAddonRef = useRef(null);

    console.log("MOunt XTermComponent");

    useLayoutEffect(() => {
        console.log("isDrawerOpen: ", isDrawerOpen);

        if (isDrawerOpen) {
            // Create a new Terminal instance
            terminalRef.current = new Terminal();

            // Attach the terminal to the container element
            terminalRef.current.open(containerRef.current);

            // Connect to the WebSocket and attach it to the terminal
            const socket = new WebSocket(socketUrl);
            const attachAddon = new AttachAddon(socket);
            const fitAddon = new FitAddon();
            terminalRef.current.loadAddon(attachAddon);
            terminalRef.current.loadAddon(fitAddon);

            // Delay the fitAddon.fit() call with setTimeout
            setTimeout(() => {
                fitAddon.fit();
                console.log('Fitting terminal');
            }, 500);

            // fitAddon.fit();
            // console.log("Fitting terminal");

            // Store the fitAddon instance in a ref so it can be accessed in the cleanup function
            fitAddonRef.current = fitAddon;

            // Add any additional logic here...
        } else {
            console.log("Closing terminal");
            // Clean up the terminal when the drawer is closed
            // terminalRef.current?.dispose();
            terminalRef.current?.dispose();
            fitAddonRef.current?.dispose();
        }

        // TODO: Close websocket connection
        // Return a cleanup function to dispose of the FitAddon instance
        return () => {
            console.log("Unmount XTermComponent");
            // terminalRef.current?.dispose();
            // fitAddonRef.current?.dispose();
        };
    }, [isDrawerOpen, socketUrl]);

    return <div style={{ "height": "100%", "width": "100%" }} ref={containerRef} />;
};

export default XTermComponent;