import { useState, useEffect } from "react";

function useLocalStorageState(key, defaultValue = '') {
  const [value, setValue] = useState(() => {
    const itemValue = localStorage.getItem(key);
    return itemValue ? JSON.parse(itemValue) : defaultValue;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue]
}

export default useLocalStorageState;
