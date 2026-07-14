export interface PlayerState {
  id: string;
  name: string;
  connected: boolean;
  totalPossessionTime: number;
  passesMade: number;
  passesRecieved: number;
}

export interface BallState {
  currentOwnerId: string;
  previousOwnerId: string;
  passCount: number;
}

export interface RoomState {
  type: "ROOM_STATE";
  roomId: string;
  gameState: "WAITING" | "RUNNING" | "FINISHED";
  maxPlayers: number;
  players: PlayerState[];
  ball: BallState;
  winnerId: string;
  winnerName: string;
  remainingSeconds: number;
}