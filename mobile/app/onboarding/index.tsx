import { View, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function OnboardingScreen() {
  return (
    <SafeAreaView className="flex-1 bg-black">
      <View className="flex-1 items-center justify-center px-6">
        <Text className="text-white text-2xl font-bold">Rate these first</Text>
        <Text className="text-gray-400 mt-2 text-center">
          We found artists you know. Rate a few to calibrate your taste.
        </Text>
      </View>
    </SafeAreaView>
  );
}
