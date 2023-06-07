import { APIGatewayProxyCallbackV2, Context } from 'aws-lambda';

export const handler = async (
  event: APIGatewayProxyCallbackV2,
  context: Context,
  callback: APIGatewayProxyCallbackV2,
): Promise<void> => {
  callback(null, event as any);
};
