# Where Are You?

A Garmin Connect IQ application to request and navigate to a friend's location.

## Concept

1.  **Request:** User A selects User B from their friend list and requests their location.
2.  **Approve:** User B receives a prompt on their watch. They can accept or deny.
3.  **Navigate:** If accepted, User A receives User B's current coordinates.
    *   The watch displays a directional arrow and distance to User B.
    *   Updates can be polled or streamed.

## Architecture

*   **Watch App:** Built with Monkey C (Garmin Connect IQ SDK).
*   **Backend:** Intermediary server to handle signaling and location relay.
