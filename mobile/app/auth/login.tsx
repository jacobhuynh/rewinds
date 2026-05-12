import { View, Text, TouchableOpacity } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { router } from "expo-router";

export default function LoginScreen() {
  return (
    <SafeAreaView className="flex-1 bg-black px-6">
      <View className="flex-1 justify-center gap-6">
        <Text className="text-white text-4xl font-bold">Rewinds</Text>
        <Text className="text-gray-400 text-lg">Your music taste, ranked.</Text>

        <TouchableOpacity
          className="bg-green-500 rounded-xl py-4 items-center mt-8"
          onPress={() => router.push("/auth/spotify-callback")}
        >
          <Text className="text-black font-bold text-base">Continue with Spotify</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="border border-gray-700 rounded-xl py-4 items-center"
          onPress={() => router.push("/auth/signup")}
        >
          <Text className="text-white font-semibold text-base">Continue with Email</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
