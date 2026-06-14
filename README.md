# Financial Data Aggregator

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Angular](https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white)](https://angular.io/)
[![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)

Aplikacja webowa służąca do agregacji danych finansowych oraz zarządzania własnym portfelem inwestycyjnym. 

---

## Stos technologiczny

### Backend
* **Język:** Go
* **Framework:** Gin
* **ORM:** GORM
* **Architektura:** Wzorzec CSR (Controllers (Handlers), Services, Repositories)

### Frontend
* **Język:** TypeScript
* **Framework:** Angular, TailwindCSS

### Infrastruktura & Bazy danych
* **Baza główna:** PostgreSQL (dane użytkowników, transakcje, portfolio)
* **Cache:** Redis (cache cen aktywów)
* **Konteneryzacja:** Docker, Docker Compose
* **Serwer WWW / Proxy:** Nginx

---

## Model architektury

System opiera się na architekturze mikroserwisowej z wyodrębnionym backendem i frontendem, które komunikują się za pomocą REST API.

```mermaid
sequenceDiagram
    participant F as Frontend (Angular)
    participant P as Proxy (Nginx)
    participant B as Backend (Go API)
    participant R as Cache (Redis)
    participant D as Baza Danych (PostgreSQL)
    participant E as Zewnętrzne API

    rect rgb(200, 220, 250)
    Note over B, E: PROCES W TLE (Worker - priceService.StartWorker)
    loop Cykliczna aktualizacja
        B->>E: Zapytanie o aktualne kursy aktywów
        E-->>B: Odpowiedź z nowymi wycenami
        B->>R: Zapisz aktualne wyceny (Cache)
    end
    end

    rect rgb(230, 250, 230)
    Note over F, R: KLIENT POBIERA KURSY (Brak autoryzacji)
    F->>P: HTTP GET /api/rates
    P->>B: Przekazanie żądania do PriceHandler
    B->>R: Odczyt zapisanych wycen z pamięci
    R-->>B: Zwrócenie błyskawicznych wyników
    B-->>P: Odpowiedź JSON z kursami
    P-->>F: Wyświetlenie danych
    end

    rect rgb(250, 230, 230)
    Note over F, D: KLIENT ZARZĄDZA PORTFELEM (Wymagany JWT)
    F->>P: HTTP POST /api/protected/portfolio
    P->>B: Przekazanie żądania (AuthMiddleware)
    B->>B: Weryfikacja tokenu JWT
    B->>D: Zapis nowego aktywa (PortfolioService -> Repo)
    D-->>B: Potwierdzenie zapisu w bazie
    B-->>P: Status 200 OK
    P-->>F: Komunikat: "item added to portfolio"
    end
```

* [Dokumentacja API](docs/models/api.md)
* [Model bazy danych](docs/models/db_model.md)

---

## Instrukcja uruchomienia (Docker)

Aplikacja jest w pełni skonteneryzowana. Aby uruchomić projekt lokalnie, upewnij się, że masz zainstalowanego **Dockera** oraz **Docker Compose**.

### Krok 1: Klonowanie repozytorium
```bash
git clone [https://github.com/r4qq/financial-data-aggregator.git](https://github.com/r4qq/financial-data-aggregator.git)
cd financial-data-aggregator
```

### Krok 2: Konfiguracja środowiska
Skopiuj przykładowy plik konfiguracyjny i dostosuj go:
```bash
cp .env.example .env
```

### Krok 3: Uruchomienie aplikacji
Aby zbudować i uruchomić wszystkie kontenery (Baza danych, Cache, Backend, Frontend) w tle, wykonaj:
```bash
docker-compose up -d --build
```

### Krok 4: Weryfikacja działania
Po udanym uruchomieniu aplikacja będzie dostępna pod następującymi adresami:
* **Frontend:** `http://localhost:80`
* **Backend API:** `http://localhost:8080/api`

---

## Live:

Projekt został wdrozony i jest dostępny pod adresem:
**[Live](https://fin.krzysztofstasiak.pl/)**

---

## Autorzy:

 * Krzysztof Stasiak
 * Stanisław Rak