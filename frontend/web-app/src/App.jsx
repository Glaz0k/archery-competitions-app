import { Route, Routes } from "react-router";
import CompetitionPage from "./components/pages/CompetitionPage";
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
        <Route path="cups/:cupId/competitions/:competitionId" element={<CompetitionPage />} />
      </Route>
    </Routes>
  );
}
