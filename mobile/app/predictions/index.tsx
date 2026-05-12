import { View, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function PredictionsScreen() {
  return (
    <SafeAreaView className="flex-1 bg-black">
      <View className="flex-1 items-center justify-center">
        <Text className="text-white text-2xl font-bold">Predictions</Text>
        <Text className="text-gray-400 mt-2">Open &amp; resolved predictions</Text>
      </View>
    </SafeAreaView>
  );
}
