import { AuthenticationError, ForbiddenError } from "apollo-server";
import { KeyRole } from "../entities/Key";
import { ArgsError } from "../lib/errors";
import { IApolloContext } from "../types";

export const isNotAnonymous = (ctx: IApolloContext) => {
  if (ctx.user.anonymous) {
    throw new AuthenticationError("Must be authenticated");
  }
};

export const isPersonalUser = (ctx: IApolloContext) => {
  isNotAnonymous(ctx);
  if (!ctx.user.key.userId || ctx.user.key.role !== KeyRole.Manage) {
    throw new ForbiddenError("Only permitted with personal login");
  }
};

export const canEditUser = (ctx: IApolloContext, userId: string) => {
  isPersonalUser(ctx);
  if (ctx.user.key.userId !== userId) {
    throw new ForbiddenError("Can only edit yourself");
  }
};

export const canReadProject = (ctx: IApolloContext, projectId: string) => {
  // TODO
};

export const canEditProject = (ctx: IApolloContext, projectId: string) => {
  // TODO
};

export const exclusiveArgs = (args: any, keys: string[]) => {
  const keysPresent = keys.map((key) => args[key] ? 1 : 0).reduce((a, b) => a + b, 0);
  if (keysPresent !== 1) {
    throw new ArgsError(`Set one and only one of these args: ${JSON.stringify(keys)}`);
  }
};
