import { Icon } from "@chakra-ui/react";
import {
  MdAddAlert,
  MdCached,
  MdDashboard,
  MdFolder,
  MdKey,
  MdPerson,
  MdTableChart,
} from "react-icons/md";

// Admin Imports
import Notifications from "views/admin/alerts-notifications";
import CachePolicies from "views/admin/caching";
import DataTables from "views/admin/dataTables";
import MainDashboard from "views/admin/default";
import ApiKeys from "views/admin/keys";
import Profile from "views/admin/profile";
import PromptManagement from "views/admin/prompt_directory";

export const routes = [
  {
    name: "Dashboard",
    layout: "/admin",
    path: "/dashboard",
    icon: <Icon as={MdDashboard} width="20px" height="20px" color="inherit" />,
    component: MainDashboard,
  },
  {
    name: "Alerts & Notifications",
    layout: "/admin",
    icon: <Icon as={MdAddAlert} width="20px" height="20px" color="inherit" />,
    path: "/alerts",
    component: Notifications,
  },
  {
    name: "Traces",
    layout: "/admin",
    icon: <Icon as={MdTableChart} width="20px" height="20px" color="inherit" />,
    path: "/data-tables",
    component: DataTables,
  },
  {
    name: "Caching",
    layout: "/admin",
    icon: <Icon as={MdCached} width="20px" height="20px" color="inherit" />,
    path: "/cost",
    component: CachePolicies,
  },

  {
    name: "Prompt Evaluation",
    layout: "/admin",
    icon: <Icon as={MdFolder} width="20px" height="20px" color="inherit" />,
    path: "/prompt-directory",
    component: PromptManagement,
  },
  {
    name: "API Keys",
    layout: "/admin",
    path: "/keys",
    icon: <Icon as={MdKey} width="20px" height="20px" color="inherit" />,
    component: ApiKeys,
  },

  {
    name: "Profile",
    layout: "/admin",
    path: "/profile",
    icon: <Icon as={MdPerson} width="20px" height="20px" color="inherit" />,
    component: Profile,
  },
];
