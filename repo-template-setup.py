import os
import glob
import re

# -- CONFIGURATION
GO_VERSION = "1.21"
TEMPLATE_REPO_OWNER = "NotCoffee418"
TEMPLATE_REPO_NAME = "GoWebsite-Boilerplate"
ACTUAL_REPO_OWNER = None
ACTUAL_REPO_NAME = None


# -- FUNCTIONS
def prompt_dev(question, default_value):
    return input(f"{question} [{default_value}]: ") or default_value


def edit_file(file_path, edit_func):
    with open(file_path, 'r') as f:
        content = f.read()

    new_content = edit_func(content)

    with open(file_path, 'w') as f:
        f.write(new_content)


def remove_gitignore_template(content):
    lines = content.split('\n')
    start_line = None
    end_line = None

    # Find start and end lines for the template
    for i, line in enumerate(lines):
        if line.strip() == "# -- TEMPLATE ONLY --":
            start_line = i
        if line.strip() == "# -- /TEMPLATE ONLY --":
            end_line = i
            break

    # Remove lines if both start and end markers are found
    if start_line is not None and end_line is not None:
        del lines[start_line:end_line+1]

    return '\n'.join(lines)


def update_mod_template_references():
    if ACTUAL_REPO_OWNER is None or ACTUAL_REPO_NAME is None:
        raise Exception("ACTUAL_REPO_OWNER and ACTUAL_REPO_NAME must be set")

    # Safety to avoid running this script from the wrong directory
    matching_files = glob.glob("repo-template-setup.py")
    if len(matching_files) == 0:
        raise FileNotFoundError("Must be run from repo directory. Script not found.")

    def replace_go_imports(content):
        pattern = f"github.com/{TEMPLATE_REPO_OWNER}/{TEMPLATE_REPO_NAME}"
        replacement = f"github.com/{ACTUAL_REPO_OWNER}/{ACTUAL_REPO_NAME}"
        return re.sub(pattern, replacement, content)

    go_files = glob.glob(f"{os.getcwd()}/**/*.go", recursive=True)
    for go_file in go_files:
        edit_file(go_file, lambda content: replace_go_imports(content))


def update_go_version_in_go_mod(content):
    lines = content.split("\n")
    new_lines = []
    for line in lines:
        if line.startswith("go "):
            new_lines.append(f"go {GO_VERSION}")
        else:
            new_lines.append(line)
    return "\n".join(new_lines)


def cleanup_repo_setup():
    os.remove(__file__)


# --- Run script
GO_VERSION = prompt_dev("Go version", GO_VERSION)
ACTUAL_REPO_OWNER = prompt_dev("Github Repository Owner", ACTUAL_REPO_OWNER)
ACTUAL_REPO_NAME = prompt_dev("Github Repository Name", ACTUAL_REPO_NAME)

print("Updating files...")
edit_file('.gitignore', remove_gitignore_template)
edit_file('go.mod', update_go_version_in_go_mod)
update_mod_template_references()
cleanup_repo_setup()
print("Commit your changes and you're good to go!")
