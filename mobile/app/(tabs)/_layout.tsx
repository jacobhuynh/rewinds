import { Tabs } from "expo-router";

export default function TabLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: "#6366f1",
        tabBarStyle: { backgroundColor: "#0f0f0f", borderTopColor: "#1f1f1f" },
        tabBarLabelStyle: { fontSize: 12 },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{ title: "Vote", tabBarIcon: () => null }}
      />
      <Tabs.Screen
        name="discover"
        options={{ title: "Discover", tabBarIcon: () => null }}
      />
      <Tabs.Screen
        name="leaderboard"
        options={{ title: "Leaderboard", tabBarIcon: () => null }}
      />
      <Tabs.Screen
        name="profile"
        options={{ title: "Profile", tabBarIcon: () => null }}
      />
    </Tabs>
  );
}
