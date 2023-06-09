import { v4 as uuidv4 } from 'uuid';
import { KinesisClient, PutRecordCommand } from '@aws-sdk/client-kinesis';
import {
  APIGatewayProxyCallbackV2,
  APIGatewayProxyEventV2,
  Context,
} from 'aws-lambda';

interface Event {
  action: 'kinesis-data-forwarder';
  data: {
    id: string;
    title: string;
  };
}

export const handler = async (
  event: APIGatewayProxyEventV2,
  _: Context,
  callback: APIGatewayProxyCallbackV2,
): Promise<void> => {
  if (!event.body) {
    callback(new Error('Missing event body'), {
      statusCode: 400,
      body: JSON.stringify({ message: 'Missing event body' }),
    });
    return;
  }

  // Tries to parse the event body as a valid Event.
  const parsedEvent = JSON.parse(event.body) as Event;
  if (parsedEvent.action === 'kinesis-data-forwarder') {
    // Setup Kinesis client.
    const client = new KinesisClient({
      signingRegion: 'us-east-1',
      endpoint: 'http://kinesis.localhost.localstack.cloud:4566',
      credentials: {
        accessKeyId: 'na',
        secretAccessKey: 'na',
      },
      region: 'us-east-1',
    });

    // Transform data to a base64 string and add an id.
    const data = {
      ...parsedEvent.data,
      id: uuidv4(),
    };
    const base64Data = Buffer.from(JSON.stringify(data));

    // Tries to send the event to Kinesis.
    const command = new PutRecordCommand({
      StreamName: 'my-kinesis-stream',
      PartitionKey: '1',
      Data: base64Data,
    });
    try {
      await client.send(command);
      callback(null, {
        statusCode: 200,
        body: JSON.stringify({
          message: 'Event received',
          encoded: base64Data,
        }),
      });
    } catch (error) {
      console.error(error);
      callback(new Error('Error sending event to Kinesis'), {
        statusCode: 500,
        body: JSON.stringify({
          message: 'Error sending event to Kinesis',
          data: error,
        }),
      });
    }

    return;
  }

  // Ignore connect/disconnect events from websocket.
  if (
    parsedEvent.action === '$connect' ||
    parsedEvent.action === '$disconnect'
  ) {
    return;
  }

  callback(new Error('Invalid event action'), {
    statusCode: 400,
    body: JSON.stringify({ message: 'Invalid event action' }),
  });
};
