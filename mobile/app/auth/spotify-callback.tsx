import { useEffect } from "react";
import { View, Text, ActivityIndicator } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function SpotifyCallbackScreen() {
  useEffect(() => {
    // Spotify OAuth callback handler — to be implemented in Sprint 1
  }, []);

  return (
    <SafeAreaView className="flex-1 bg-black">
      <View className="flex-1 items-center justify-center gap-4">
        <ActivityIndicator color="#6366f1" size="large" />
        <Text className="text-gray-400">Connecting Spotify…</Text>
      </View>
    </SafeAreaView>
  );
}
