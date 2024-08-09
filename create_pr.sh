#!/bin/bash

# Get GITHUB_REPO_URL from environment variables
GITHUB_REPO_URL=${GITHUB_REPO_URL}

# Check if GITHUB_REPO_URL is set
if [ -z "$GITHUB_REPO_URL" ]; then
  echo "GITHUB_REPO_URL is not set in the environment variables."
  exit 1
fi

# Clone the repository
git clone "$GITHUB_REPO_URL" repo
cd repo || exit

# Calculate and store the week range from Monday to Friday
week_range=$(echo "$(date -d "monday this week" +"%d-%b-%Y")-to-$(date -d "friday this week" +"%d-%b-%Y")")

# Create a new branch for your changes with the current date
BRANCH_NAME="WeeklyUpdate-$week_range"
git checkout -b "$BRANCH_NAME"

# Make your changes here
# For example, let's create a new file
echo "Some changes" > changes.txt

# Add and commit the changes
git add .
git commit -m "Weekly update for $week_range"

# Push the changes to the new branch
git push origin "$BRANCH_NAME"

# Create a pull request using GitHub CLI
gh pr create --title "Weekly Updates for $week_range" --body "This PR includes Weeklyupdates" --base master --head "$BRANCH_NAME"

echo "Pull request created successfully."