import type { RoomState } from "../types/game";
export const socket = new WebSocket("ws://localhost:8080/ws");


socket.onopen = () => {
  console.log("🟢 Connected to Go WebSocket server");
};

socket.onmessage = (event) => {
  console.log("📨 Message from backend:", event.data);
};

socket.onerror = (error) => {
  console.error("🔴 WebSocket error:", error);
};

socket.onclose = () => {
  console.log("⚪ WebSocket connection closed");
};
export function joinPlayer(
  playerId: string,
  playerName: string
) {
  if (socket.readyState !== WebSocket.OPEN) {
    console.error("🔴 WebSocket is not connected");
    return;
  }

  const event = {
    type: "PLAYER_JOINED",
    playerId: playerId,
    playerName: playerName,
  };

  socket.send(JSON.stringify(event));

  console.log("📤 Player join event sent:", event);
}
export function subscribeToRoomState(
  callback: (roomState: RoomState) => void
) {
  const handleMessage = (event: MessageEvent) => {
    console.log("📨 Message from backend:", event.data);

    try {
      const data = JSON.parse(event.data);

      if (data.type === "ROOM_STATE") {
        callback(data as RoomState);
      }
    } catch (error) {
      console.error(
        "🔴 Failed to parse WebSocket message:",
        error
      );
    }
  };

  socket.addEventListener("message", handleMessage);

  return () => {
    socket.removeEventListener(
      "message",
      handleMessage
    );
  };
}
export function passBall(
  playerId: string,
  targetPlayerId: string
) {
  if (socket.readyState !== WebSocket.OPEN) {
    console.error("🔴 WebSocket is not connected");
    return;
  }

  const event = {
    type: "BALL_PASSED",
    playerId,
    targetPlayerId,
  };

  socket.send(JSON.stringify(event));

  console.log("🏓 Ball pass event sent:", event);
}
export function configureGame(
  gameDurationSeconds: number
) {
  if (socket.readyState !== WebSocket.OPEN) {
    console.error("🔴 WebSocket is not connected");
    return;
  }

  const event = {
    type: "GAME_CONFIGURED",
    gameDurationSeconds,
  };

  socket.send(JSON.stringify(event));

  console.log(
    "⚙️ Game configuration sent:",
    event
  );
}