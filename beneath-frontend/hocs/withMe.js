import { Query } from "react-apollo";

import { TokenConsumer } from "./auth";
import { QUERY_ME } from "../queries/user";

const withMe = (Component) => {
  return (props) => (
    <TokenConsumer>
      {({ token }) => {
        if (token) {
          return (
            <Query query={QUERY_ME}>
              {({ loading, error, data }) => {
                if (data) {
                  let { me } = data;
                  return <Component {...props} me={me} />
                } else if (error) {
                  console.log("withMe error: ", error);
                }
                return null;
              }}
            </Query>
          );
        } else {
          return (
            <Component {...props} me={null} />
          )
        }
      }}
    </TokenConsumer>
  );
}

export default withMe;
