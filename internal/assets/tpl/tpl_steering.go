//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Foundation steering file names.
const (
	// SteeringNameProduct is the name for the product context file.
	SteeringNameProduct = "product"
	// SteeringNameTech is the name for the technology stack file.
	SteeringNameTech = "tech"
	// SteeringNameStructure is the name for the project structure file.
	SteeringNameStructure = "structure"
	// SteeringNameWorkflow is the name for the development workflow file.
	SteeringNameWorkflow = "workflow"
)

// Foundation steering file descriptions.
const (
	// SteeringDescProduct describes the product context file.
	SteeringDescProduct = "Product context, goals, and target users"
	// SteeringDescTech describes the technology stack file.
	SteeringDescTech = "Technology stack, constraints, " +
		"and dependencies"
	// SteeringDescStructure describes the project structure file.
	SteeringDescStructure = "Project structure and " +
		"directory conventions"
	// SteeringDescWorkflow describes the development workflow file.
	SteeringDescWorkflow = "Development workflow and process rules"
)

// Foundation steering file body templates.
const (
	// SteeringBodyProduct is the body for the product context file.
	SteeringBodyProduct = "# Product Context\n\n" +
		"Describe the product, its goals, " +
		"and target users.\n"
	// SteeringBodyTech is the body for the technology stack file.
	SteeringBodyTech = "# Technology Stack\n\n" +
		"Describe the technology stack, " +
		"constraints, and key dependencies.\n"
	// SteeringBodyStructure is the body for the project structure
	// file.
	SteeringBodyStructure = "# Project Structure\n\n" +
		"Describe the project layout " +
		"and directory conventions.\n"
	// SteeringBodyWorkflow is the body for the development workflow
	// file.
	SteeringBodyWorkflow = "# Development Workflow\n\n" +
		"Describe the development workflow, " +
		"branching strategy, and process rules.\n"
)
