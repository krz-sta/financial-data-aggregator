# Specyfikacja API

Base URL: http://localhost:8080

Content-Type: application/json

Autoryzacja (gdzie wymagana): Nagłówek HTTP `Authorization: Bearer <token>`


## 1. Moduł: Uwierzytelnianie (/api/auth)

### POST /api/auth/register

Opis: Rejestracja nowego użytkownika.

Body:

{
  "email": "user@example.com",
  "displayName": "User1", 
  "password": "securepassword123" 
}

Statusy:

 - 201 Created - Zwraca obiekt z danymi użytkownika pod kluczem "data"

 - 400 Bad Request - Błąd walidacji wejścia (wymagane: email, min 4 znaki dla displayName, min 8 znaków dla password)

 - 500 Internal Server Error - Błąd po stronie serwera/bazy


### POST /api/auth/login

Opis: Logowanie i generowanie JWT.

Body:

{
  "email": "user@example.com",
  "password": "securepassword123"
}

Statusy:

 - 200 OK - Zwraca token
    {
    "token": "eyJhbGciOiJIUzI1NiIsInR5c..."
    }

 - 400 Bad Request - Błąd walidacji
 
 - 401 Unauthorized - Błędne dane logowania (Invalid email or password)

## 2. Moduł: Profil i Portfel (Wymaga JWT - prefiks: /api/protected)

### POST /api/protected/profile

Opis: Pobranie danych profilu oraz przypisanego portfela dla zalogowanego użytkownika na podstawie tokenu (ID z contextu).

Body: Brak

Statusy:
 - 200 OK - Zwraca profil wraz z zawartym portfelem
 - 404 Not Found - Nie znaleziono użytkownika
 - 500 Internal Server Error - Brak ID użytkownika w kontekście

### POST /api/protected/portfolio

Opis: Dodanie nowego aktywa do portfela.

Body:
{
  "symbol": "BTC",
  "amount": 0.5
}

Statusy:

- 200 OK - Sukces
 {
   "message": "item added to portoflio"
 }

 - 400 Bad Request - Błąd formatu wejścia
 
 - 401 Unauthorized - Brak autoryzacji
 
 - 500 Internal Server Error - Błąd serwera przy dodawaniu

### DELETE /api/protected/portfolio/:id

Opis: Usunięcie elementu z portfela (gdzie :id to UUID usuwanej pozycji).

Statusy:
 - 200 OK - Sukces
 {
   "message": "item deleted from portoflio"
 }

 - 401 Unauthorized - Brak autoryzacji
 - 500 Internal Server Error - Błąd bazy/usunięcia

## 3. Moduł: Aktywa i Kursy (/api/assets, /api/rates)

### GET /api/assets

Opis: Pobiera listę wspieranych aktywów do śledzenia.

Statusy:

 - 200 OK - Zwraca tablicę obsługiwanych aktywów i ich metadanych (bez id api).

### GET /api/rates

Opis: Pobiera aktualne kursy aktywów (korzysta z cache Redis).

Statusy:
 - 200 OK - Zwraca mapę kursów

### GET /api/rates/history/:symbol

Opis: Pobiera historię cen dla danego symbolu waluty/kryptowaluty.

Statusy:

 - 200 OK - Zwraca tablicę z historycznymi punktami (Timestamp, Price)
 - 500 Internal Server Error - Błąd usługi historii

## 4. Moduł: Diagnostyka (/api/health)

### GET /api/health/db

Opis: Sprawdza połączenie z główną bazą PostgreSQL.

Statusy:
 - 200 OK - { "status": "UP" }
 - 500 Internal Server Error - { "error": "DB is down" } lub błąd bazy

### GET /api/health/redis

Opis: Sprawdza połączenie z pamięcią podręczną Redis (metoda Ping).

Statusy:
 - 200 OK - { "status": "UP" }
 - 500 Internal Server Error - { "error": "Redis is down" }