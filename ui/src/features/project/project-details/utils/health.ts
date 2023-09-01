import unknownIcon from '@ui/assets/icons/heart-exclamation.svg';
import healthyIcon from '@ui/assets/icons/heart-green.svg';
import unhealthyIcon from '@ui/assets/icons/heart-red.svg';
import { HealthState } from '@ui/gen/v1alpha1/types_pb';

export const healthStateToString = (status?: HealthState): string => {
  switch (status) {
    case HealthState.HEALTHY:
      return 'Healthy';
    case HealthState.UNHEALTHY:
      return 'Unhealthy';
    case HealthState.UNKNOWN:
      return 'Unknown';
    default:
      return '';
  }
};

export const healthStateToIcon = (status?: HealthState) => {
  switch (status) {
    case HealthState.HEALTHY:
      return healthyIcon;
    case HealthState.UNHEALTHY:
      return unhealthyIcon;
    case HealthState.UNKNOWN:
      return unknownIcon;
    default:
      return '';
  }
};