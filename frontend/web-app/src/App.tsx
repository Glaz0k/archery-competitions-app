import { Route, Routes } from "react-router";
import { ProtectedRoute } from "./features";
import CompetitorsContent from "./pages/competition/CompetitorsContent";
import GroupsContent from "./pages/competition/GroupsContent";
import CompetitionPage from "./pages/CompetitionPage";
import CompetitorsPage from "./pages/CompetitorsPage";
import ContentLayout from "./pages/ContentLayout";
import CupPage from "./pages/CupPage";
import CupsPage from "./pages/CupsPage";
import IndividualGroupPage from "./pages/IndividualGroupPage";
import NotFoundPage from "./pages/NotFoundPage";
import SignInPage from "./pages/SignInPage";

export default function App() {
  return (
    <Routes>
      <Route path="*" element={<NotFoundPage />} />
      <Route path="sign-in" element={<SignInPage />} />
      <Route element={<ProtectedRoute />}>
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
                <Route path="individual-groups/:groupId" element={<IndividualGroupPage />}>
                  <Route path=":groupSection" element={null} />
                </Route>
              </Route>
            </Route>
          </Route>
          <Route path="competitors" element={<CompetitorsPage />} />
        </Route>
      </Route>
    </Routes>
  );
}
