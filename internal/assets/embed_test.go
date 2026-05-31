//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

// Embedded-asset tests are split by concern across this package:
//
//   - templates_test.go — context-file templates (TestGetTemplate,
//     TestListTemplates)
//   - project_test.go    — project-root files (TestClaudeMd,
//     TestProjectFile, TestMakefileCtx)
//   - skills_test.go     — Claude skills + references (TestListSkills,
//     TestSkillContent, TestSkillReference, …)
//   - why_test.go        — why-docs (TestWhyDoc, TestListWhyDocs)
//   - plugin_test.go     — plugin manifest (TestPluginVersion)
//   - schema_test.go     — .ctxrc JSON schema (TestSchema,
//     TestSchemaCoversCtxRC)
//   - hooks_test.go      — hook message registry (TestHookMessageRegistry,
//     TestListHookMessages, TestHookMessage_ReadVariant)
//
// Two tests that exercise read/ subpackages live in those packages
// instead, to avoid an assets → read/X → assets import cycle:
//
//   - TestDescKeysResolve  → read/desc/desc_test.go (needs lookup.Init)
//   - default-permissions  → read/lookup/perm_test.go (TestPermAllowListDefault,
//     TestPermDenyListDefault — need Init)
