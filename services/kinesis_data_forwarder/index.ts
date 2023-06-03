import {
  APIGatewayProxyEventV2,
  APIGatewayProxyCallbackV2,
  Context,
} from 'aws-lambda';

export const handler = async (
  event: APIGatewayProxyEventV2,
  context: Context,
  callback: APIGatewayProxyCallbackV2,
): Promise<void> => {
  const b = JSON.stringify(event.body);
  callback(null, {
    statusCode: 200,
    body: JSON.stringify({
      message: `Hello ${b}`,
    }),
  });
};
