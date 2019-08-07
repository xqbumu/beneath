import React from "react";
import isEmail from "validator/lib/isEmail";
import { Mutation } from "react-apollo";

import Button from "@material-ui/core/Button";
import DeleteIcon from "@material-ui/icons/Delete";
import Grid from "@material-ui/core/Grid";
import IconButton from "@material-ui/core/IconButton";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemAvatar from "@material-ui/core/ListItemAvatar";
import ListItemText from "@material-ui/core/ListItemText";
import ListItemSecondaryAction from "@material-ui/core/ListItemSecondaryAction";
import TextField from "@material-ui/core/TextField";
import { makeStyles } from "@material-ui/core";

import Avatar from "../Avatar";
import Loading from "../Loading";
import NextMuiLink from "../NextMuiLink";
import VSpace from "../VSpace";

import { QUERY_PROJECT, ADD_MEMBER, REMOVE_MEMBER } from "../../apollo/queries/project";

const useStyles = makeStyles((theme) => ({
  addMemberContainer: {
    [theme.breakpoints.up("md")]: {
      width: theme.breakpoints.values.sm,
    },
  },
}));

const ManageMembers = ({ project, editable }) => {
  return (
    <React.Fragment>
      <ViewMembers project={project} editable={editable} />
      {editable && <VSpace units={2} />}
      {editable && <AddMember project={project} />}
    </React.Fragment>
  );
};

export default ManageMembers;

const ViewMembers = ({ project, editable }) => {
  const classes = useStyles();
  return (
    <List>
      {project.users.map(({ userID, username, name, photoURL }) => (
        <ListItem
          key={userID}
          component={NextMuiLink} as={`/users/${userID}`} href={`/user?id=${userID}`}
          disableGutters button
        >
          <ListItemAvatar><Avatar size="list" label={name} src={photoURL} /></ListItemAvatar>
          <ListItemText primary={name} secondary={username} />
          {editable && (project.users.length > 1) && (
            <ListItemSecondaryAction>
              <Mutation mutation={REMOVE_MEMBER} update={(cache, { data: { removeUserFromProject } }) => {
                const projectName = project.name;
                if (removeUserFromProject) {
                  const { projectByName } = cache.readQuery({ query: QUERY_PROJECT, variables: { name: projectName } });
                  cache.writeQuery({
                    query: QUERY_PROJECT,
                    variables: { name: projectName },
                    data: { projectByName: { ...projectByName, users: projectByName.users.filter((user) => user.userID !== userID) } },
                  });
                }
              }}>
                {(removeUserFromProject, { loading, error }) => (
                  <IconButton edge="end" aria-label="Delete" onClick={() => {
                    removeUserFromProject({ variables: { projectID: project.projectID, userID: userID } });
                  }}>
                    {loading ? <Loading size={20} /> : <DeleteIcon />}
                  </IconButton>
                )}
              </Mutation>
            </ListItemSecondaryAction>
          )}
        </ListItem>
      ))}
    </List>
  );
};

const AddMember = ({ project }) => {
  const [email, setEmail] = React.useState("");
  const classes = useStyles();
  return (
    <Mutation mutation={ADD_MEMBER} update={(cache, { data: { addUserToProject } }) => {
      const user = addUserToProject;
      const query = cache.readQuery({ query: QUERY_PROJECT, variables: { name: project.name } });
      cache.writeQuery({
        query: QUERY_PROJECT,
        variables: { name: project.name },
        data: { projectByName: { ...query.projectByName, users: query.projectByName.users.concat([user]) } },
      });
    }} onCompleted={() => setEmail("")}>
      {(addUserToProject, { loading, error }) => (
        <form onSubmit={(e) => {
          e.preventDefault();
          addUserToProject({ variables: { email, projectID: project.projectID } });
        }}>
          <Grid container alignItems={"center"} spacing={2} className={classes.addMemberContainer}>
            <Grid item xs={true}>
              <TextField id="email" type="email" label="Add Member" placeholder="Email" value={email}
                fullWidth disabled={loading} error={!!error} helperText={error && error.graphQLErrors[0].message}
                onChange={(event) => setEmail(event.target.value)}
              />
            </Grid>
            <Grid item>
              <Button type="submit" variant="outlined" color="primary" disabled={loading || !isEmail(email)}>
                Add member
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
    </Mutation>
  );
};