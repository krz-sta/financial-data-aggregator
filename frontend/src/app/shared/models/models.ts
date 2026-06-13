export interface Asset {
  symbol: string;
  name: string;
  type: string;
}

export interface PortfolioItem {
  id: string;
  symbol: string;
  amount: number;
}

export interface UserProfile {
  id: string;
  email: string;
  displayName: string;
  portfolio: PortfolioItem[];
  createdAt: string;
}

export interface RatesResponse {
  [symbol: string]: number;
}

export interface HistoryPoint {
  timestamp: number;
  price: number;
}

export interface FiatRates {
  [currency: string]: number;
}
