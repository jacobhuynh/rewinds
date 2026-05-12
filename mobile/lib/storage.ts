import * as SecureStore from "expo-secure-store";

const JWT_KEY = "rewinds_jwt";

export async function saveToken(token: string) {
  await SecureStore.setItemAsync(JWT_KEY, token);
}

export async function getToken(): Promise<string | null> {
  return SecureStore.getItemAsync(JWT_KEY);
}

export async function deleteToken() {
  await SecureStore.deleteItemAsync(JWT_KEY);
}
