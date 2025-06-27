import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Pet, CartItem } from '../types';

interface CartContextType {
  cartItems: CartItem[];
  cartCount: number;
  addToCart: (pet: Pet) => void;
  removeFromCart: (petId: string) => void;
  clearCart: () => void;
  isInCart: (petId: string) => boolean;
  getCartPetIds: () => string[];
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};

interface CartProviderProps {
  children: ReactNode;
}

export const CartProvider: React.FC<CartProviderProps> = ({ children }) => {
  const [cartItems, setCartItems] = useState<CartItem[]>([]);

  useEffect(() => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      setCartItems(JSON.parse(savedCart));
    }
  }, []);

  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(cartItems));
  }, [cartItems]);

  const addToCart = (pet: Pet) => {
    setCartItems(prev => {
      const exists = prev.find(item => item.pet.id === pet.id);
      if (!exists) {
        return [...prev, { pet, addedAt: Date.now() }];
      }
      return prev;
    });
  };

  const removeFromCart = (petId: string) => {
    setCartItems(prev => prev.filter(item => item.pet.id !== petId));
  };

  const clearCart = () => {
    setCartItems([]);
  };

  const isInCart = (petId: string) => {
    return cartItems.some(item => item.pet.id === petId);
  };

  const getCartPetIds = () => {
    return cartItems.map(item => item.pet.id);
  };

  const cartCount = cartItems.length;

  return (
    <CartContext.Provider
      value={{
        cartItems,
        cartCount,
        addToCart,
        removeFromCart,
        clearCart,
        isInCart,
        getCartPetIds,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};