import { useEffect, useState } from "react";

import {
  configureGame,
  joinPlayer,
  passBall,
  subscribeToRoomState,
} from "./websocket/socket";

import type { RoomState } from "./types/game";

function App() {
  // Latest room state from Go backend
  const [roomState, setRoomState] =
    useState<RoomState | null>(null);

  // Join form values
  const [playerId, setPlayerId] = useState("");
  const [playerName, setPlayerName] = useState("");

  // Player identity for this browser
  const [joinedPlayerId, setJoinedPlayerId] =
    useState("");

  // Local visual countdown
  // Resynchronized whenever backend sends ROOM_STATE
  const [countdown, setCountdown] = useState(0);
  const [gameDurationSeconds, setGameDurationSeconds] =
  useState(30);

  // Subscribe to live ROOM_STATE messages
  useEffect(() => {
    const unsubscribe = subscribeToRoomState(
      (newRoomState) => {
        console.log(
          "🏓 React received room state:",
          newRoomState
        );

        // Store latest backend room state
        setRoomState(newRoomState);

        // Synchronize countdown with backend
        setCountdown(
          newRoomState.remainingSeconds
        );
      }
    );

    return unsubscribe;
  }, []);

  // Local countdown:
  // 10 -> 9 -> 8 -> ... -> 0
  useEffect(() => {
    // Countdown should run only
    // while game is RUNNING
    if (roomState?.gameState !== "RUNNING") {
      return;
    }

    // Stop when countdown reaches zero
    if (countdown <= 0) {
      return;
    }

    const timer = window.setTimeout(() => {
      setCountdown((currentCountdown) =>
        Math.max(currentCountdown - 1, 0)
      );
    }, 1000);

    // Cleanup old timer
    return () => {
      window.clearTimeout(timer);
    };
  }, [
    countdown,
    roomState?.gameState,
  ]);
  

  // Join game
  const handleJoin = () => {
    const trimmedPlayerId = playerId.trim();
    const trimmedPlayerName = playerName.trim();

    if (
      trimmedPlayerId === "" ||
      trimmedPlayerName === ""
    ) {
      console.error(
        "🔴 Player ID and Player Name are required"
      );

      return;
    }

    // Send PLAYER_JOINED event
    // to Go WebSocket backend
    joinPlayer(
      trimmedPlayerId,
      trimmedPlayerName
    );

    // Remember which player belongs
    // to this browser
    setJoinedPlayerId(trimmedPlayerId);
  };

  // Check whether this browser's
  // player currently owns the ball
  const isCurrentOwner =
    roomState?.ball.currentOwnerId ===
    joinedPlayerId;

  // Current 2-player demo:
  // find the other player
  const targetPlayer =
    roomState?.players.find(
      (player) =>
        player.id !== joinedPlayerId
    );

  return (
    <main className="game-app">
      <div className="game-container">

        {/* HEADER */}
        <header className="game-header">
          <div>
            <h1 className="game-title">
              🏓 Real-Time Ping Pong
            </h1>

            <p className="game-subtitle">
              Event-driven multiplayer game
              powered by Go + WebSockets
            </p>
          </div>

          <div className="status-badge">
            {roomState?.gameState ?? "OFFLINE"}
          </div>
        </header>
                {/* GAME CONFIGURATION */}
        <section className="panel join-panel">
          <select
            value={gameDurationSeconds}
            onChange={(event) =>
              setGameDurationSeconds(
                Number(event.target.value)
              )
            }
          >
            <option value={30}>30 Seconds</option>
            <option value={60}>1 Minute</option>
            <option value={120}>2 Minutes</option>
            <option value={300}>5 Minutes</option>
          </select>

          <button
            className="primary-button"
            onClick={() =>
              configureGame(gameDurationSeconds)
            }
          >
            ⚙️ Apply Game Duration
          </button>
        </section>

        {/* JOIN PANEL */}
        <section className="panel join-panel">
          <input
            type="text"
            placeholder="Player ID"
            value={playerId}
            onChange={(event) =>
              setPlayerId(event.target.value)
            }
          />

          <input
            type="text"
            placeholder="Player Name"
            value={playerName}
            onChange={(event) =>
              setPlayerName(event.target.value)
            }
          />

          <button
            className="primary-button"
            onClick={handleJoin}
          >
            Join Game
          </button>
        </section>

        {/* WAITING FOR FIRST ROOM STATE */}
        {roomState === null ? (
          <section className="panel empty-state">
            <h2>Waiting for players</h2>

            <p>
              Join the room to receive live
              game state from the Go server.
            </p>
          </section>
        ) : (
          <>
            {/* COUNTDOWN */}
            {roomState.gameState === "RUNNING" && (
              <section className="countdown-card">
                <span className="countdown-value">
                  {countdown}
                </span>

                <span className="countdown-label">
                  seconds to auto-pass
                </span>
              </section>
            )}

            {/* GAME ARENA */}
            <section className="panel arena">

              {/* PLAYER CARDS */}
              <div className="players-grid">
                {roomState.players.map(
                  (player) => {
                    const ownsBall =
                      roomState.ball
                        .currentOwnerId ===
                      player.id;

                    return (
                      <article
                        key={player.id}
                        className={
                          ownsBall
                            ? "player-card owner"
                            : "player-card"
                        }
                      >
                        <h2 className="player-name">
                          {player.name}
                        </h2>

                        {/* BALL OWNER LABEL */}
                        {ownsBall && (
                          <p className="owner-label">
                            🏓 HOLDING THE BALL
                          </p>
                        )}

                        {/* PLAYER STATS */}
                        <div className="stats-grid">

                          {/* POSSESSION */}
                          <div className="stat">
                            <span className="stat-value">
                              {
                                player.totalPossessionTime
                              }
                              s
                            </span>

                            <span className="stat-label">
                              Possession
                            </span>
                          </div>

                          {/* PASSES MADE */}
                          <div className="stat">
                            <span className="stat-value">
                              {player.passesMade}
                            </span>

                            <span className="stat-label">
                              Passes Made
                            </span>
                          </div>

                          {/* PASSES RECEIVED */}
                          <div className="stat">
                            <span className="stat-value">
                              {
                                player.passesRecieved
                              }
                            </span>

                            <span className="stat-label">
                              Received
                            </span>
                          </div>
                        </div>
                      </article>
                    );
                  }
                )}
              </div>

              {/* ANIMATED BALL TRACK */}
              <div className="ball-track">
                <div
                  className={`ball-indicator ${
                    roomState.ball.currentOwnerId ===
                    roomState.players[0]?.id
                      ? "ball-left"
                      : "ball-right"
                  }`}
                >
                  🏓
                </div>
              </div>
            </section>

            {/* GAME CONTROLS */}
            <section className="panel game-footer">

              {/* ROOM INFO */}
              <div>
                <strong>
                  Room {roomState.roomId}
                </strong>

                <p>
                  Players:{" "}
                  {roomState.players.length}
                  {" / "}
                  {roomState.maxPlayers}

                  {" • "}

                  Total Passes:{" "}
                  {roomState.ball.passCount}
                </p>
              </div>

              {/* PASS BUTTON */}
              {roomState.gameState ===
                "RUNNING" &&
                isCurrentOwner &&
                targetPlayer && (
                  <button
                    className="pass-button"
                    onClick={() =>
                      passBall(
                        joinedPlayerId,
                        targetPlayer.id
                      )
                    }
                  >
                    🏓 Pass to{" "}
                    {targetPlayer.name}
                  </button>
                )}
            </section>

            {/* WINNER */}
            {roomState.gameState ===
              "FINISHED" && (
              <section className="winner-card">
                <h2>
                  🏆{" "}
                  {roomState.winnerName}
                  {" "}
                  Wins!
                </h2>

                <p>
                  Holding the ball at the
                  final moment.
                </p>
              </section>
            )}
          </>
        )}
      </div>
    </main>
  );
}

export default App;