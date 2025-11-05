optimism_package = import_module("github.com/tokamak-network/optimism-package/main.star")

def run(plan, args):
    # Upload local artifacts before running optimism-package
    l1_artifacts = plan.upload_files(
        src = "./artifacts/l1-artifacts",
        name = "l1-artifacts"
    )

    l2_artifacts = plan.upload_files(
        src = "./artifacts/l2-artifacts",
        name = "l2-artifacts"
    )

    # Load local simple.yaml config
    # just delegate to optimism-package with local args
    optimism_package.run(plan, args)
