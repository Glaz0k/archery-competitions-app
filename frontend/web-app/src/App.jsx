import { Route, Routes } from "react-router";
import CompetitorsContent from "./components/pages/competition/CompetitorsContent";
import GroupsContent from "./components/pages/competition/GroupsContent";
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
        <Route path="cups/:cupId" element={<CompetitionPage />}>
          <Route path="competitions/:competitionId" element={<GroupsContent />} />
          <Route path="competitions/:competitionId/competitors" element={<CompetitorsContent />} />
        </Route>
      </Route>
    </Routes>
  );
}
