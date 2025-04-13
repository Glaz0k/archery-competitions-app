import { Route, Routes } from "react-router";
import ContentLayout from "./components/pages/ContentLayout";
import CupPage from "./components/pages/CupPage";
import CupsPage from "./components/pages/CupsPage";
import LoginPage from "./components/pages/LoginPage";

export default function App() {
  return (
    <Routes>
      <Route path="login" element={<LoginPage />} />
      <Route element={<ContentLayout />}>
        <Route path="cups" element={<CupsPage />} />
        <Route path="cups/:cupId" element={<CupPage />} />
      </Route>
    </Routes>
  );
}
