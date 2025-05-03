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
        <Route path="cups">
          <Route index element={<CupsPage />} />
          <Route path=":cupId">
            <Route index element={<CupPage />} />
            <Route path="competitions/:competitionId">
              <Route element={<CompetitionPage />}>
                <Route index element={<GroupsContent />} />
                <Route path="competitors" element={<CompetitorsContent />} />
              </Route>
              {/*<Route path="individual-groups/:groupId" element={<IndividualGroupPage />}>
                <Route path=":groupSection" element={null} />
              </Route>*/}
            </Route>
          </Route>
        </Route>
      </Route>
    </Routes>
  );
}
