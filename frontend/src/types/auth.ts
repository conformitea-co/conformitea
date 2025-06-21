export interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
  provider: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
}

export interface AuthStore extends AuthState {
  setUser: (user: User | null) => void;
  clearUser: () => void;
}
