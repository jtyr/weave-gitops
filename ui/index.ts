import AutomationsTable from "./components/AutomationsTable";
import BucketDetail from "./components/BucketDetail";
import Button from "./components/Button";
import DataTable, { SortType } from "./components/DataTable";
import FilterableTable, {
  filterConfig,
  FilterConfigCallback,
  filterByStatusCallback,
  filterByTypeCallback,
} from "./components/FilterableTable";
import Flex from "./components/Flex";
import FluxRuntime from "./components/FluxRuntime";
import Footer from "./components/Footer";
import GithubDeviceAuthModal from "./components/GithubDeviceAuthModal";
import GitRepositoryDetail from "./components/GitRepositoryDetail";
import HelmChartDetail from "./components/HelmChartDetail";
import HelmReleaseDetail from "./components/HelmReleaseDetail";
import HelmRepositoryDetail from "./components/HelmRepositoryDetail";
import Icon, { IconType } from "./components/Icon";
import Interval from "./components/Interval";
import KubeStatusIndicator from "./components/KubeStatusIndicator";
import KustomizationDetail from "./components/KustomizationDetail";
import Link from "./components/Link";
import LoadingPage from "./components/LoadingPage";
import Page from "./components/Page";
import RepoInputWithAuth from "./components/RepoInputWithAuth";
import SourceLink from "./components/SourceLink";
import SourcesTable from "./components/SourcesTable";
import Timestamp from "./components/Timestamp";
import UserSettings from "./components/UserSettings";
import AppContextProvider from "./contexts/AppContext";
import CoreClientContextProvider, {
  UnAuthorizedInterceptor,
} from "./contexts/CoreClientContext";
import AuthContextProvider, { Auth, AuthCheck } from "./contexts/AuthContext";
import CallbackStateContextProvider from "./contexts/CallbackStateContext";
import {
  Automation,
  useGetHelmRelease,
  useGetKustomization,
  useListAutomations,
} from "./hooks/automations";
import { useFeatureFlags } from "./hooks/featureflags";
import { useListFluxRuntimeObjects } from "./hooks/flux";
import { useIsAuthenticated } from "./hooks/gitprovider";
import { useListSources } from "./hooks/sources";
import { Applications as applicationsClient } from "./lib/api/applications/applications.pb";
import { Core as coreClient } from "./lib/api/core/core.pb";
import { FluxObjectKind } from "./lib/api/core/types.pb";
import {
  clearCallbackState,
  getCallbackState,
  getProviderToken,
} from "./lib/storage";
import { muiTheme, theme } from "./lib/theme";
import { V2Routes } from "./lib/types";
import { statusSortHelper } from "./lib/utils";
import OAuthCallback from "./pages/OAuthCallback";
import SignIn from "./pages/SignIn";
import SubRouterTabs, { RouterTab } from "./components/SubRouterTabs";
import InfoList from "./components/InfoList";
import EventsTable from "./components/EventsTable";
import { useGetObject } from "./hooks/objects";
import Metadata from "./components/Metadata";
import { fluxObjectKindToKind } from "./lib/objects";
import { formatURL } from "./lib/nav";

export {
  AppContextProvider,
  applicationsClient,
  Auth,
  AuthCheck,
  AuthContextProvider,
  Automation,
  AutomationsTable,
  BucketDetail,
  Button,
  CallbackStateContextProvider,
  clearCallbackState,
  coreClient,
  CoreClientContextProvider,
  DataTable,
  EventsTable,
  Flex,
  FilterableTable,
  filterByStatusCallback,
  filterByTypeCallback,
  filterConfig,
  FilterConfigCallback,
  FluxObjectKind,
  fluxObjectKindToKind,
  FluxRuntime,
  Footer,
  formatURL,
  getCallbackState,
  getProviderToken,
  GithubDeviceAuthModal,
  GitRepositoryDetail,
  HelmChartDetail,
  HelmReleaseDetail,
  HelmRepositoryDetail,
  Icon,
  IconType,
  InfoList,
  Interval,
  KubeStatusIndicator,
  KustomizationDetail,
  Link,
  LoadingPage,
  Metadata,
  muiTheme,
  OAuthCallback,
  Page,
  RepoInputWithAuth,
  RouterTab,
  SignIn,
  SortType,
  SourceLink,
  SourcesTable,
  statusSortHelper,
  SubRouterTabs,
  theme,
  Timestamp,
  UnAuthorizedInterceptor,
  useFeatureFlags,
  useGetHelmRelease,
  useGetKustomization,
  useGetObject,
  useIsAuthenticated,
  useListAutomations,
  useListFluxRuntimeObjects,
  useListSources,
  UserSettings,
  V2Routes,
};
