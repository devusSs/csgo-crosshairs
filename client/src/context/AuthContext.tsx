import React , { createContext, useState } from "react";


const AuthContext = createContext({});

export const AuthProvider = ({ children }: any) => {
  const [auth, setAuth] = useState({'role' : localStorage.getItem('role') || ''});
  
  
  return (
    <AuthContext.Provider value={{ auth, setAuth }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
