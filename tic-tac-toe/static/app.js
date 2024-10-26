document.addEventListener('DOMContentLoaded', () => {
    let board = ["", "", "", "", "", "", "", "", ""];
    let currentPlayer = 'X';
    let isGameActive = true;
    let player1Name = "";
    let player2Name = "";

    const startGameButton = document.getElementById('start-game');
    const cells = document.querySelectorAll('.cell');
    const result = document.getElementById('result');
    const gameBoard = document.getElementById('game-board');
    const player1Input = document.getElementById('player1');
    const player2Input = document.getElementById('player2');
    const messageBox = document.getElementById('message-box'); // New

    const winningConditions = [
        [0, 1, 2], [3, 4, 5], [6, 7, 8],
        [0, 3, 6], [1, 4, 7], [2, 5, 8],
        [0, 4, 8], [2, 4, 6]
    ];

    function checkWinner() {
        let roundWon = false;
        for (let i = 0; i < winningConditions.length; i++) {
            const [a, b, c] = winningConditions[i];
            if (board[a] && board[a] === board[b] && board[a] === board[c]) {
                roundWon = true;
                break;
            }
        }
        if (roundWon) {
            const winner = currentPlayer === 'X' ? player1Name : player2Name;
            declareWinner(winner);
            saveGameResult(player1Name, player2Name, winner); // Save the game result
            isGameActive = false;
            return;
        }
        if (!board.includes("")) {
            declareTie();
        }
    }

    function declareWinner(winner) {
        result.textContent = `Congratulations! ${winner} wins!`;
    }

    function declareTie() {
        result.textContent = `It's a tie!`;
    }

    function handleCellClick(event) {
        const clickedCellIndex = event.target.getAttribute('data-index');
        if (board[clickedCellIndex] !== "" || !isGameActive) return;

        board[clickedCellIndex] = currentPlayer;
        event.target.textContent = currentPlayer;

        if (currentPlayer === 'X') {
            event.target.classList.add('player-x');
        } else {
            event.target.classList.add('player-o');
        }

        checkWinner();
        currentPlayer = currentPlayer === 'X' ? 'O' : 'X';
    }

    function startGame() {
        player1Name = player1Input.value;
        player2Name = player2Input.value;
        if (player1Name === "" || player2Name === "") {
            alert("Please enter both player names.");
            return;
        }

        board = ["", "", "", "", "", "", "", "", ""];
        cells.forEach(cell => {
            cell.textContent = "";
            cell.classList.remove('player-x', 'player-o');
        });
        result.textContent = "";
        messageBox.textContent = ""; // Clear message
        currentPlayer = 'X';
        isGameActive = true;

        gameBoard.classList.remove('hidden');
    }

    function saveGameResult(player1, player2, winner) {
        const gameResult = {
            player1_name: player1,
            player2_name: player2,
            winner_name: winner
        };

        fetch('/save-result', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(gameResult),
        })
        .then(response => response.json())
        .then(data => {
            messageBox.textContent = data.message; // Display success message
        })
        .catch(error => {
            messageBox.textContent = 'Error saving winner to DB';
        });
    }

    startGameButton.addEventListener('click', startGame);
    cells.forEach(cell => cell.addEventListener('click', handleCellClick));
});
