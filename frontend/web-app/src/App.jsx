import { Routes, Route } from "react-router";
import LoginPage from "./components/pages/LoginPage";
import ContentLayout from "./components/pages/ContentLayout";
import CupsPage from "./components/pages/CupsPage";

export default function App() {
  return (
    <Routes>
      <Route path="login" element={<LoginPage />} />
      <Route element={<ContentLayout />}>
        <Route path="cups" element={<CupsPage />} />
      </Route>
    </Routes>
  );
}
