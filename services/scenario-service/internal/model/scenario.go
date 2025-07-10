// # структура сценариев и шагов
package model

type Scenario struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
}

type Step struct {
	ID         int64   `json:"id"`
	ScenarioID int64   `json:"scenario_id"`
	Text       string  `json:"text"`
	Place      string  `json:"place"`
	Time       string  `json:"time"`
	Price      float64 `json:"price"`
}

// CREATE TABLE scenarios (
//     id SERIAL PRIMARY KEY,
//     title TEXT NOT NULL,
//     location TEXT,
//     price NUMERIC
// );

// CREATE TABLE steps (
//     id SERIAL PRIMARY KEY,
//     scenario_id INTEGER REFERENCES scenarios(id) ON DELETE CASCADE,
//     title TEXT,
//     description TEXT,
//     time TEXT
// );
