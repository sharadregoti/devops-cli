export interface ResourceTypes {
  id: string;
  name: string;
  data: Datum[];
  specificActions: SpecificAction[];
}

interface SpecificAction {
  name: string;
  key_binding: string;
  output_type: string;
  args?: any;
  schema?: any;
  execution: Execution;
}

interface Execution {
  cmd: string;
  is_long_running: boolean;
  user_input: Userinput;
  server_input: Serverinput;
}

interface Serverinput {
  required: boolean;
}

interface Userinput {
  required: boolean;
  args?: any;
}

interface Datum {
  data: string[];
  color: string;
}