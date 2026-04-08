import { Tabs } from "expo-router";

const CompanyLayout = () => {
    return (
        <Tabs screenOptions={{ 
            headerShown: false,
        }} >
            <Tabs.Screen name="index" options={{ 
                title: "Your companies"
            }} />
            <Tabs.Screen name="account" options={{ 
                title: "Your account"
            }} />
        </Tabs>
    );
}

export default CompanyLayout;