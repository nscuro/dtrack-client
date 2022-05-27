// Package notification provides the functionality to process notifications sent by Dependency-Track.
//
// This package contains partially redundant struct definitions,
// because notification content differs from their respective API
// representations in a few ways.
//
// Dependency-Track has special serialization logic for notifications,
// which is defined here: https://github.com/DependencyTrack/dependency-track/blob/4.5.0/src/main/java/org/dependencytrack/util/NotificationUtil.java
package notification
