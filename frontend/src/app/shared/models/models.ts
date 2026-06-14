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

export type RatesResponse = Record<string, number>;

export interface HistoryPoint {
  timestamp: number;
  price: number;
}

export type FiatRates = Record<string, number>;
