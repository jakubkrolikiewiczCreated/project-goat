package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var (
	completed = make(map[int]bool)
	mu        sync.Mutex
)



func taskPage(w http.ResponseWriter, r *http.Request) {
	etStr := r.URL.Query().Get("et")


	taskNum, err := strconv.Atoi(etStr)
	if err != nil || taskNum < 1 || taskNum > 12 {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

    taskTitle, taskText := getTask(taskNum)


html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
    <title>Zadanie %d</title>

    <style>
        body {
            font-family: Arial, sans-serif;
            background: #f5f7fa;
            margin: 0;
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
        }

        .card {
            background: white;
            width: 500px;
            max-width: 90%%;
            padding: 40px;
            border-radius: 16px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            text-align: center;
        }

        .task-number {
            color: #666;
            font-size: 14px;
            margin-bottom: 10px;
        }

        h1 {
            color: #222;
        }

        p {
            color: #555;
        }

        button {
            border: none;
            background: #2563eb;
            color: white;
            padding: 14px 28px;
            font-size: 16px;
            border-radius: 8px;
            cursor: pointer;
        }

        button.done {
            background: #16a34a;
        }

        .message {
            margin-top: 15px;
            color: #16a34a;
            display: none;
            font-weight: bold;
        }
    </style>
</head>
<body>

<div class="card">

    <div class="task-number">Zadanie %d</div>

    <h1>%s</h1>

    <p>%s</p>

    <button id="doneBtn" onclick="markDone()">Ukończone</button>

    <div id="message" class="message">✓ Quest zakończony</div>
</div>

<script>
const taskNumber = %d;

function markDone() {
    fetch('/task-done', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ task: taskNumber })
    })
    .then(r => r.json())
    .then(() => {
        document.getElementById('doneBtn').innerText = 'Completed';
        document.getElementById('doneBtn').classList.add('done');
        document.getElementById('doneBtn').disabled = true;
        document.getElementById('message').style.display = 'block';
    });
}
</script>

</body>
</html>
`,
taskNum,
taskNum,
taskTitle,     
taskText,    
taskNum)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func taskDone(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Task int `json:"task"`
	}

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	completed[req.Task] = true
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func completedCount(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count := len(completed)
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}

func main() {
	http.HandleFunc("/task", taskPage)
	http.HandleFunc("/task-done", taskDone)
	http.HandleFunc("/completed-count", completedCount)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func getTask(taskNumber int) (string, string) {
    taskTexts := map[int]string{
	1:  "Jesteś uczestnikiem MasterChef Junior - przygotuj dowolnego drinka (w liczbie 3), a następnie przekaż go do spróbowania Jury - zadanie zostanie zaliczone, gdy 2 sędziom zasmakuje ten drink.",
	2:  "Odbij 10 razy dowolną częścią ciała piłkę (masz na to 5 prób).",
	3:  "Przytocz z pamięci 5 kultowych cytatów z serii gier FIFA (od Szpakowskiego, Laskowskiego lub zagranicznych komentatorów).",
	4:  "Wygłoś oryginalny, rymowany wiersz-toast, który porwie wszystkich do wzniesienia kieliszków.",
	5:  "Rozegraj partię w klasyczną grę karcianą „Wojna” i zakończ pojedynek, wygrywając z przeciwnikiem przewagą dokładnie 3 punktów.",
	6:  "Wymień z nazwiska 5 czarnoskórych piłkarzy, którzy w historii reprezentowali barwy Widzewa Łódź.",
	7:  "Wykaż się pamięcią i wymień 10 nauczycieli z liceum SLO wraz z przedmiotami, których nauczają lub nauczali.",
	8:  "Wymień 7 Wielkich Polaków, którzy zmienili bieg historii, nauki lub kultury, i krótko uzasadnij swój wybór. Zadanie zostanie zaliczone, gdy wszyscy sędziowie będą zgodni.",
	9:  "Znajdź Remka i odpalcie wspólnego, relaksacyjnego dymka cygara w klimatycznym miejscu.",
	10: "Zorganizuj epickiego grilla z Bartkiem - zabezpiecz prowiant, rozpal ogień i przypieczcie coś dobrego.",
	11: "Zasiądź do stołu z Kubą i rozegrajcie emocjonującą partię w kości jak w Wiedźminie 2.",
	12: "Stocz pojedynek 1v1 w wymienianiu postaci z Summoner's Rift.",
	13: "Zjedz plasterek cytryny bez zmrużenia oka lub wykonaj karne zadanie alkoholowe/smakowe jako finałowy bonus spotkania.",
}

 taskTitles := map[int]string{
	1:  "Przepierdoliłem kupę siana",
	2:  "Piłkarski GOAT",
	3:  "Głos Komentatora",
	4:  "Mistrz Ceremonii",
	5:  "Strateg Wojny",
	6:  "Widzewska Afryka",
	7:  "Kronika SLO",
	8:  "Panteon Chwały",
	9:  "Dymne Przymierze",
	10: "Władca Rusztu",
	11: "Kości zostały rzucone",
	12: "Arena Bohaterów",
	13: "Kwaśny Finał (Bonus)",
}

    return taskTitles[taskNumber], taskTexts[taskNumber]
}