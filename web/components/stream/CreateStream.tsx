import { useMutation, useQuery } from "@apollo/client";
import { Typography } from "@material-ui/core";
import { Field, Form, Formik } from "formik";
import { useRouter } from "next/router";
import React, { FC } from "react";

import { QUERY_PROJECTS_FOR_USER } from "../../apollo/queries/project";
import { STAGE_STREAM } from "../../apollo/queries/stream";
import { StreamSchemaKind } from "../../apollo/types/globalTypes";
import { ProjectsForUser, ProjectsForUserVariables } from "../../apollo/types/ProjectsForUser";
import { StageStream, StageStreamVariables } from "../../apollo/types/StageStream";
import useMe from "../../hooks/useMe";
import { toURLName, toBackendName } from "../../lib/names";
import { handleSubmitMutation, SelectField as FormikSelectField, TextField as FormikTextField } from "../formik";
import SubmitControl from "../forms/SubmitControl";

interface Project {
  organization: { name: string };
  name: string;
  displayName?: string;
}

export interface CreateProjectProps {
  preselectedProject?: Project;
}

const CreateStream: FC<CreateProjectProps> = ({ preselectedProject }) => {
  const me = useMe();
  const router = useRouter();
  const [stageStream] = useMutation<StageStream, StageStreamVariables>(STAGE_STREAM, {
    onCompleted: (data) => {
      if (data?.stageStream) {
        const orgName = toURLName(data.stageStream.project.organization.name);
        const projName = toURLName(data.stageStream.project.name);
        const streamName = toURLName(data.stageStream.name);
        const href = `/stream?organization_name=${orgName}&project_name=${projName}&stream_name=${streamName}`;
        const as = `/${orgName}/${projName}/${streamName}`;
        router.replace(href, as, { shallow: true });
      }
    },
  });

  const { data, loading, error } = useQuery<ProjectsForUser, ProjectsForUserVariables>(QUERY_PROJECTS_FOR_USER, {
    variables: { userID: me?.personalUserID || "" },
    skip: !me,
  });

  const initialValues = {
    project: data?.projectsForUser?.length ? (preselectedProject ? preselectedProject : data.projectsForUser[0]) : null,
    name: "",
    schemaKind: StreamSchemaKind.GraphQL,
    schema: "",
  };

  return (
    <Formik
      initialStatus={error?.message}
      initialValues={initialValues}
      onSubmit={async (values, actions) =>
        handleSubmitMutation(
          values,
          actions,
          stageStream({
            variables: {
              organizationName: toBackendName(values.project?.organization.name || ""),
              projectName: toBackendName(values.project?.name || ""),
              streamName: toBackendName(values.name),
              schemaKind: values.schemaKind,
              schema: values.schema,
            },
          })
        )
      }
    >
      {({ isSubmitting, status }) => (
        <Form>
          <Typography component="h2" variant="h1" gutterBottom>
            Create stream
          </Typography>
          <Field
            name="project"
            validate={(proj?: Project) => {
              if (!proj) {
                return "Select a project for the stream";
              }
            }}
            component={FormikSelectField}
            label="Project"
            required
            loading={loading}
            options={data?.projectsForUser || []}
            getOptionLabel={(option: Project) => `${option.organization.name}/${option.name}`}
            getOptionSelected={(option: Project, value: Project) => {
              return option.name === value.name;
            }}
          />
          <Field
            name="name"
            validate={(val: string) => {
              if (!val || val.length < 3 || val.length > 40) {
                return "Stream names should be between 3 and 40 characters long";
              }
              if (!val.match(/^[_\-a-z][_\-a-z0-9]+$/)) {
                return "Stream names should consist of lowercase letters, numbers, underscores and dashes (cannot start with a number)";
              }
            }}
            component={FormikTextField}
            label="Name"
            required
          />
          <Field
            name="schemaKind"
            validate={(kind?: string) => {
              if (!kind) {
                return "Select a schema kind for the stream";
              }
            }}
            component={FormikSelectField}
            label="Schema language"
            required
            disableClearable
            options={[StreamSchemaKind.GraphQL]}
          />
          <Field
            name="schema"
            validate={(schema?: string) => {
              if (!schema) {
                return "You must provide a valid schema";
              }
            }}
            component={FormikTextField}
            label="Schema"
            required
            multiline
            rows={10}
            rowsMax={200}
          />
          <SubmitControl label="Create stream" errorAlert={status} disabled={isSubmitting} />
        </Form>
      )}
    </Formik>
  );
};

export default CreateStream;
