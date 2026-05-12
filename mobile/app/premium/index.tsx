import { View, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function PremiumScreen() {
  return (
    <SafeAreaView className="flex-1 bg-black">
      <View className="flex-1 items-center justify-center">
        <Text className="text-white text-2xl font-bold">Go Premium</Text>
        <Text className="text-gray-400 mt-2">Higher bet limits, custom themes, and more</Text>
      </View>
    </SafeAreaView>
  );
}
