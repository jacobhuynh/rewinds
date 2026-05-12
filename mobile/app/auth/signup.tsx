import { View, Text, TextInput, TouchableOpacity } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function SignupScreen() {
  return (
    <SafeAreaView className="flex-1 bg-black px-6">
      <View className="flex-1 justify-center gap-4">
        <Text className="text-white text-3xl font-bold mb-4">Create account</Text>

        <TextInput
          placeholder="Email"
          placeholderTextColor="#6b7280"
          className="bg-gray-900 text-white rounded-xl px-4 py-4 border border-gray-800"
          keyboardType="email-address"
          autoCapitalize="none"
        />
        <TextInput
          placeholder="Username"
          placeholderTextColor="#6b7280"
          className="bg-gray-900 text-white rounded-xl px-4 py-4 border border-gray-800"
          autoCapitalize="none"
        />
        <TextInput
          placeholder="Password"
          placeholderTextColor="#6b7280"
          className="bg-gray-900 text-white rounded-xl px-4 py-4 border border-gray-800"
          secureTextEntry
        />

        <TouchableOpacity className="bg-indigo-500 rounded-xl py-4 items-center mt-2">
          <Text className="text-white font-bold text-base">Sign up</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
