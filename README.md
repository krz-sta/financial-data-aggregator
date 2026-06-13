# Financial Data Aggregator

System agregacji i monitorowania portfela aktywów finansowych. Zrealizowany w architekturze rozproszonej, podzielony na niezależne komponenty backendu i frontendu, uruchamiany w środowisku kontenerowym.

## 1. Architektura i Stos Technologiczny

Backend:
- Język: Go 1.25
- Framework HTTP: Gin Gonic
- ORM: GORM (driver: postgres)
- Uwierzytelnianie: JWT (JSON Web Tokens) bezstanowe
- Testy: testify, testcontainers-go

Frontend:
- Framework: Angular 21 (TypeScript)
- Arkusze stylów: SCSS, TailwindCSS
- Serwer HTTP (PROD): Nginx (konfiguracja: nginx.conf)

Bazy danych:
- Główna baza danych: PostgreSQL 15 (przechowywanie profili użytkowników, portfeli)
- Cache: Redis 9 (buforowanie zapytań rynkowych, ograniczanie rate limitów API zewnętrznych)

Infrastruktura:
- Konteneryzacja: Docker
- Orkiestracja: Docker Compose (osobne pliki dla DEV i PROD)
- CI/CD: GitHub Actions (automatyzacja testów, linterów i budowy obrazów)

## 2. Zgodność z The Twelve-Factor App

1. Codebase: Jedno repozytorium (monorepo) z kodem aplikacji i infrastruktury.
2. Dependencies: Zależności izolowane przez Dockera oraz pliki go.mod i package.json.
3. Config: Konfiguracja oddzielona od kodu, oparta na zmiennych środowiskowych (plik .env).
4. Backing services: PostgreSQL i Redis uruchamiane jako odrębne kontenery i podłączane przez sieć Dockera.
5. Stateless: Backend nie przechowuje stanu sesji w pamięci (zastosowanie JWT).
6. Build, release, run: Wielofazowe budowanie obrazów (Multi-stage builds) oddzielające kompilację od środowiska uruchomieniowego (obrazy oparte na Alpine).

## 3. Zmienne Środowiskowe (.env)

Zmienne wymagane do uruchomienia środowiska:

DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
REDIS_HOST
REDIS_PORT
REDIS_PASSWORD
ROUTER_HOST
ROUTER_PORT
JWT_SECRET

## 4. Instrukcja Uruchomienia

Wymagania: Docker, Docker Compose.
Przed uruchomieniem należy skopiować plik .env.example do .env i uzupełnić parametry.

Środowisko deweloperskie (Hot Reloading):
> docker compose up --build
- Frontend: http://localhost:4200
- Backend: http://localhost:8080

Środowisko produkcyjne (Frontend serwowany przez Nginx, minimalne obrazy):
> docker compose -f docker-compose.prod.yml up --build -d
- Frontend (Nginx): http://localhost:80
- Backend: http://localhost:8080

## 5. Konfiguracja Testów Integracyjnych (Testcontainers)

Testy warstwy repozytorium zaimplementowano z użyciem biblioteki testcontainers-go. 
Proces testowania "go test ./...":
1. Pobiera obraz postgres:15-alpine z rejestru.
2. Uruchamia efemeryczny kontener testowy.
3. Wykonuje Auto-Migracje GORM.
4. Przeprowadza operacje DML sprawdzające zgodność zapytań z dialektem PostgreSQL.
5. Wysyła sygnał przerwania (Terminate) i usuwa kontener z hosta.

## 6. Wykonali
 - Krzysztof Stasiak
 - Stanisław Rak