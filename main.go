package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type WeatherData struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
		Humidity  int     `json:"humidity"`
		FeelsLike float64 `json:"feelslike_c"`
		WindKph   float64 `json:"wind_kph"`
	} `json:"current"`
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Aplikasi Cuaca</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            max-width: 500px;
            width: 100%;
        }

        h1 {
            color: #667eea;
            text-align: center;
            margin-bottom: 30px;
            font-size: 2em;
        }

        .search-box {
            display: flex;
            gap: 10px;
            margin-bottom: 30px;
        }

        input {
            flex: 1;
            padding: 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 16px;
            transition: border 0.3s;
        }

        input:focus {
            outline: none;
            border-color: #667eea;
        }

        button {
            padding: 15px 30px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 10px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: transform 0.2s;
        }

        button:hover {
            transform: translateY(-2px);
        }

        .weather-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 15px;
            padding: 30px;
            color: white;
            text-align: center;
        }

        .location {
            font-size: 1.5em;
            margin-bottom: 20px;
            font-weight: 600;
        }

        .weather-icon {
            width: 100px;
            height: 100px;
            margin: 0 auto;
        }

        .temperature {
            font-size: 4em;
            font-weight: 700;
            margin: 20px 0;
        }

        .condition {
            font-size: 1.3em;
            margin-bottom: 20px;
            opacity: 0.9;
        }

        .details {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 20px;
            margin-top: 20px;
            padding-top: 20px;
            border-top: 1px solid rgba(255, 255, 255, 0.3);
        }

        .detail-item {
            text-align: center;
        }

        .detail-label {
            font-size: 0.9em;
            opacity: 0.8;
            margin-bottom: 5px;
        }

        .detail-value {
            font-size: 1.2em;
            font-weight: 600;
        }

        .error {
            background: #ff6b6b;
            color: white;
            padding: 15px;
            border-radius: 10px;
            text-align: center;
            margin-bottom: 20px;
        }

        .loading {
            text-align: center;
            color: #667eea;
            font-size: 1.2em;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌤️ Cuaca Hari Ini</h1>
        
        <form action="/" method="GET" class="search-box">
            <input type="text" name="city" placeholder="Masukkan nama kota..." value="{{.Query}}" required>
            <button type="submit">Cari</button>
        </form>

        {{if .Error}}
            <div class="error">{{.Error}}</div>
        {{end}}

        {{if .Weather}}
        <div class="weather-card">
            <div class="location">
                {{.Weather.Location.Name}}, {{.Weather.Location.Country}}
            </div>
            
            <img src="https:{{.Weather.Current.Condition.Icon}}" alt="Weather icon" class="weather-icon">
            
            <div class="temperature">
                {{printf "%.0f" .Weather.Current.TempC}}°C
            </div>
            
            <div class="condition">
                {{.Weather.Current.Condition.Text}}
            </div>
            
            <div class="details">
                <div class="detail-item">
                    <div class="detail-label">Terasa Seperti</div>
                    <div class="detail-value">{{printf "%.0f" .Weather.Current.FeelsLike}}°C</div>
                </div>
                <div class="detail-item">
                    <div class="detail-label">Kelembaban</div>
                    <div class="detail-value">{{.Weather.Current.Humidity}}%</div>
                </div>
                <div class="detail-item">
                    <div class="detail-label">Angin</div>
                    <div class="detail-value">{{printf "%.0f" .Weather.Current.WindKph}} km/h</div>
                </div>
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>
`

type PageData struct {
	Weather *WeatherData
	Error   string
	Query   string
}

func getWeather(city string) (*WeatherData, error) {
	apiKey := "6c0e6b47bbf74fc38fd81842250610" // https://www.weatherapi.com/signup.aspx
	
	baseURL := "http://api.weatherapi.com/v1/current.json"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", city)
	params.Add("aqi", "no")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("kota tidak ditemukan atau API key tidak valid")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherData
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("weather").Parse(htmlTemplate))
	
	data := PageData{
		Query: r.URL.Query().Get("city"),
	}

	if data.Query != "" {
		weather, err := getWeather(data.Query)
		if err != nil {
			data.Error = err.Error()
		} else {
			data.Weather = weather
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	
	fmt.Println("🌤️  Server berjalan di http://localhost:8080")
	fmt.Println("Tekan Ctrl+C untuk berhenti")
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}