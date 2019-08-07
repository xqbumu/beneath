import { ApolloClient } from "apollo-boost";
import React, { FC } from "react";
import { Query, withApollo } from "react-apollo";
import { SubscriptionClient } from "subscriptions-transport-ws";

import { QUERY_LATEST_RECORDS } from "../../apollo/queries/local/records";
import { LatestRecords, LatestRecordsVariables } from "../../apollo/types/LatestRecords";
import { QueryStream } from "../../apollo/types/QueryStream";
import { GATEWAY_URL_WS } from "../../lib/connection";
import RecordsTable from "../RecordsTable";
import { Schema } from "./schema";

type StreamLatestProps = QueryStream & { client: ApolloClient<any> };

interface StreamLatestState {
  error: string | undefined;
}

class StreamLatest extends React.Component<StreamLatestProps, StreamLatestState> {
  private apollo: ApolloClient<any>
  private subscription: SubscriptionClient | undefined;
  private schema: Schema;

  constructor(props: StreamLatestProps) {
    super(props);
    this.apollo = props.client;
    this.schema = new Schema(props.stream, true);
    this.state = {
      error: undefined,
    };
  }

  public componentWillUnmount() {
    if (this.subscription) {
      this.subscription.close(true);
    }
  }

  public componentDidMount() {
    if (this.subscription) {
      return;
    }

    const self = this;

    this.subscription = new SubscriptionClient(`${GATEWAY_URL_WS}/ws`, {
      reconnect: true,
    });

    const request = {
      query: this.props.stream.currentStreamInstanceID || undefined,
    };

    const apolloVariables = {
      projectName: this.props.stream.project.name,
      streamName: this.props.stream.name,
    };

    this.subscription.request(request).subscribe({
      next: (result) => {
        const queryData = self.apollo.cache.readQuery({
          query: QUERY_LATEST_RECORDS,
          variables: apolloVariables,
        }) as any;

        if (!queryData.latestRecords) {
          queryData.latestRecords = [];
        }

        const recordID = self.schema.makeUniqueIdentifier(result.data);

        queryData.latestRecords = queryData.latestRecords.filter((record: any) => record.recordID !== recordID);

        queryData.latestRecords.unshift({
          __typename: "Record",
          recordID,
          data: result.data,
          sequenceNumber: result.data && result.data["@meta"] && result.data["@meta"].sequence_number,
        });

        self.apollo.writeQuery({
          query: QUERY_LATEST_RECORDS,
          variables: apolloVariables,
          data: queryData,
        });
      },
      error: (error) => {
        self.setState({
          error: error.message,
        });
      },
      complete: () => {
        if (!self.state.error) {
          self.setState({
            error: "Unexpected completion of subscription",
          });
        }
        self.subscription = undefined;
      },
    });
  }

  public render() {
    const variables = {
      projectName: this.props.stream.project.name,
      streamName: this.props.stream.name,
      limit: 100,
    };

    return (
      <Query<LatestRecords, LatestRecordsVariables>
        query={QUERY_LATEST_RECORDS}
        variables={variables}
        fetchPolicy="cache-and-network"
      >
        {({ loading, error, data }) => {
          const errorMsg = error || this.state.error;
          if (errorMsg) {
            return <p>Error: {JSON.stringify(error)}</p>;
          }

          loading = loading || !!this.subscription;

          return (
            <RecordsTable schema={this.schema} loading={loading} records={data ? data.latestRecords : null} />
          );
        }}
      </Query>
    );
  }
}

export default withApollo(StreamLatest);
