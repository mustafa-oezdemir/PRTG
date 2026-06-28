export interface ManualApiMethod {
  label: string;
  value: string;
  description: string;
}

export const manualApiMethods: ManualApiMethod[] = [
  {
    label: 'Get Sensor Details',
    value: 'getsensordetails.json',
    description: 'Get detailed information about a sensor',
  },
  {
    label: 'Get Status',
    value: 'getstatus.htm',
    description: 'Retrieve system status information',
  },
];
