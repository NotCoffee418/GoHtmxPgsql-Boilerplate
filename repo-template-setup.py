import os
import glob
import re

# -- CONFIGURATION
GO_VERSION = "1.21"
TEMPLATE_REPO_OWNER = "NotCoffee418"
TEMPLATE_REPO_NAME = "GoHtmxPgsql-Boilerplate"


# -- FUNCTIONS
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
    actual_repo_owner = os.environ['GITHUB_REPOSITORY'].split('/')[0]
    actual_repo_name = os.environ['GITHUB_REPOSITORY'].split('/')[1]

    def replace_go_imports(content):
        pattern = f"github.com/{TEMPLATE_REPO_OWNER}/{TEMPLATE_REPO_NAME}"
        replacement = f"github.com/{actual_repo_owner}/{actual_repo_name}"
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
    os.remove(".github/workflows/repo-template-setup.yml")


# --- Run script
edit_file('.gitignore', remove_gitignore_template)
edit_file('go.mod', update_go_version_in_go_mod)
update_mod_template_references()
cleanup_repo_setup()
