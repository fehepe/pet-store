import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface AuthContextType {
  isAuthenticated: boolean;
  storeId: string | null;
  customerName: string | null;
  login: (username: string, password: string, storeId: string) => boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [storeId, setStoreId] = useState<string | null>(null);
  const [customerName, setCustomerName] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    const savedStoreId = localStorage.getItem('storeId');
    const savedCustomerName = localStorage.getItem('customerName');
    
    if (token && savedStoreId && savedCustomerName) {
      setIsAuthenticated(true);
      setStoreId(savedStoreId);
      setCustomerName(savedCustomerName);
    }
  }, []);

  const login = (username: string, password: string, storeIdParam: string) => {
    // Create basic auth token
    const token = btoa(`${username}:${password}`);
    localStorage.setItem('authToken', token);
    localStorage.setItem('storeId', storeIdParam);
    localStorage.setItem('customerName', username);
    
    setIsAuthenticated(true);
    setStoreId(storeIdParam);
    setCustomerName(username);
    
    return true;
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('storeId');
    localStorage.removeItem('customerName');
    setIsAuthenticated(false);
    setStoreId(null);
    setCustomerName(null);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, storeId, customerName, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};