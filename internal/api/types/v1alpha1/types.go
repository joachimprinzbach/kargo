package v1alpha1

import (
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	kubemetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	kargoapi "github.com/akuity/kargo/api/v1alpha1"
	typesmetav1 "github.com/akuity/kargo/internal/api/types/metav1"
	"github.com/akuity/kargo/internal/version"
	svcv1alpha1 "github.com/akuity/kargo/pkg/api/service/v1alpha1"
	"github.com/akuity/kargo/pkg/api/v1alpha1"
)

func FromProjectProto(p *svcv1alpha1.Project) *unstructured.Unstructured {
	if p == nil {
		return nil
	}
	u := &unstructured.Unstructured{}
	u.SetAPIVersion(kargoapi.GroupVersion.String())
	u.SetKind("Project")
	u.SetCreationTimestamp(kubemetav1.NewTime(p.GetCreateTime().AsTime()))
	u.SetName(p.GetName())
	return u
}

func FromStageProto(s *v1alpha1.Stage) *kargoapi.Stage {
	if s == nil {
		return nil
	}
	var status kargoapi.StageStatus
	if s.GetStatus() != nil {
		status = *FromStageStatusProto(s.GetStatus())
	}
	var objectMeta kubemetav1.ObjectMeta
	if s.GetMetadata() != nil {
		objectMeta = *typesmetav1.FromObjectMetaProto(s.GetMetadata())
	}
	return &kargoapi.Stage{
		TypeMeta: kubemetav1.TypeMeta{
			APIVersion: kargoapi.GroupVersion.String(),
			Kind:       "Stage",
		},
		ObjectMeta: objectMeta,
		Spec:       FromStageSpecProto(s.GetSpec()),
		Status:     status,
	}
}

func FromStageSpecProto(s *v1alpha1.StageSpec) *kargoapi.StageSpec {
	return &kargoapi.StageSpec{
		Subscriptions:       FromSubscriptionsProto(s.GetSubscriptions()),
		PromotionMechanisms: FromPromotionMechanismsProto(s.GetPromotionMechanisms()),
	}
}

func FromStageStatusProto(s *v1alpha1.StageStatus) *kargoapi.StageStatus {
	if s == nil {
		return nil
	}
	history := make(kargoapi.SimpleFreightStack, len(s.GetHistory()))
	for idx, freight := range s.GetHistory() {
		history[idx] = *FromSimpleFreightProto(freight)
	}
	return &kargoapi.StageStatus{
		CurrentFreight: FromSimpleFreightProto(s.GetCurrentFreight()),
		History:        history,
		Health:         FromHealthProto(s.GetHealth()),
		Error:          s.GetError(),
	}
}

func FromFreightProto(f *v1alpha1.Freight) *kargoapi.Freight {
	if f == nil {
		return nil
	}
	var objectMeta kubemetav1.ObjectMeta
	if f.GetMetadata() != nil {
		objectMeta = *typesmetav1.FromObjectMetaProto(f.GetMetadata())
	}
	commits := make([]kargoapi.GitCommit, len(f.GetCommits()))
	for idx, commit := range f.GetCommits() {
		commits[idx] = *FromGitCommitProto(commit)
	}
	images := make([]kargoapi.Image, len(f.GetImages()))
	for idx, image := range f.GetImages() {
		images[idx] = *FromImageProto(image)
	}
	charts := make([]kargoapi.Chart, len(f.GetCharts()))
	for idx, chart := range f.GetCharts() {
		charts[idx] = *FromChartProto(chart)
	}
	qualifications :=
		make(map[string]kargoapi.Qualification, len(f.Status.Qualifications))
	for stageName := range f.Status.Qualifications {
		qualifications[stageName] = kargoapi.Qualification{}
	}
	return &kargoapi.Freight{
		TypeMeta: kubemetav1.TypeMeta{
			APIVersion: kargoapi.GroupVersion.String(),
			Kind:       "Freight",
		},
		ObjectMeta: objectMeta,
		ID:         f.GetId(),
		Commits:    commits,
		Images:     images,
		Charts:     charts,
		Status: kargoapi.FreightStatus{
			Qualifications: qualifications,
		},
	}
}

func FromSimpleFreightProto(s *v1alpha1.SimpleFreight) *kargoapi.SimpleFreight {
	if s == nil {
		return nil
	}
	commits := make([]kargoapi.GitCommit, len(s.GetCommits()))
	for idx, commit := range s.GetCommits() {
		commits[idx] = *FromGitCommitProto(commit)
	}
	images := make([]kargoapi.Image, len(s.GetImages()))
	for idx, image := range s.GetImages() {
		images[idx] = *FromImageProto(image)
	}
	charts := make([]kargoapi.Chart, len(s.GetCharts()))
	for idx, chart := range s.GetCharts() {
		charts[idx] = *FromChartProto(chart)
	}
	return &kargoapi.SimpleFreight{
		ID:      s.GetId(),
		Commits: commits,
		Images:  images,
		Charts:  charts,
	}
}

func FromWarehouseProto(w *v1alpha1.Warehouse) *kargoapi.Warehouse {
	if w == nil {
		return nil
	}
	var objectMeta kubemetav1.ObjectMeta
	if w.GetMetadata() != nil {
		objectMeta = *typesmetav1.FromObjectMetaProto(w.GetMetadata())
	}
	var status kargoapi.WarehouseStatus
	if w.GetStatus() != nil {
		status = *FromWarehouseStatusProto(w.GetStatus())
	}
	return &kargoapi.Warehouse{
		TypeMeta: kubemetav1.TypeMeta{
			APIVersion: kargoapi.GroupVersion.String(),
			Kind:       "Warehouse",
		},
		ObjectMeta: objectMeta,
		Spec:       FromWarehouseSpecProto(w.GetSpec()),
		Status:     status,
	}
}

func FromWarehouseSpecProto(s *v1alpha1.WarehouseSpec) *kargoapi.WarehouseSpec {
	if s == nil {
		return nil
	}
	subscriptions := make([]kargoapi.RepoSubscription, 0, len(s.GetSubscriptions()))
	for _, subscription := range s.GetSubscriptions() {
		if subscription == nil {
			continue
		}
		subscriptions = append(subscriptions, *FromRepoSubscriptionProto(subscription))
	}
	return &kargoapi.WarehouseSpec{
		Subscriptions: subscriptions,
	}
}

func FromWarehouseStatusProto(s *v1alpha1.WarehouseStatus) *kargoapi.WarehouseStatus {
	if s == nil {
		return nil
	}
	return &kargoapi.WarehouseStatus{}
}

func FromGitCommitProto(g *v1alpha1.GitCommit) *kargoapi.GitCommit {
	if g == nil {
		return nil
	}
	return &kargoapi.GitCommit{
		RepoURL:           g.GetRepoUrl(),
		ID:                g.GetId(),
		Branch:            g.GetBranch(),
		HealthCheckCommit: g.GetHealthCheckCommit(),
	}
}

func FromImageProto(i *v1alpha1.Image) *kargoapi.Image {
	if i == nil {
		return nil
	}
	return &kargoapi.Image{
		RepoURL: i.GetRepoUrl(),
		Tag:     i.GetTag(),
	}
}

func FromChartProto(c *v1alpha1.Chart) *kargoapi.Chart {
	if c == nil {
		return nil
	}
	return &kargoapi.Chart{
		RegistryURL: c.GetRegistryUrl(),
		Name:        c.GetName(),
		Version:     c.GetVersion(),
	}
}

func FromHealthProto(h *v1alpha1.Health) *kargoapi.Health {
	if h == nil {
		return nil
	}
	argocdAppStates := make([]kargoapi.ArgoCDAppStatus, len(h.GetArgocdApps()))
	for i, argocdAppState := range h.GetArgocdApps() {
		argocdAppStates[i] = FromArgoCDAppStateProto(argocdAppState)
	}
	return &kargoapi.Health{
		Status:     kargoapi.HealthState(h.GetStatus()),
		Issues:     h.GetIssues(),
		ArgoCDApps: argocdAppStates,
	}
}

func FromArgoCDAppStateProto(
	a *v1alpha1.ArgoCDAppState,
) kargoapi.ArgoCDAppStatus {
	return kargoapi.ArgoCDAppStatus{
		Namespace:    a.GetNamespace(),
		Name:         a.GetName(),
		HealthStatus: FromArgoCDAppHealthStatusProto(a.GetHealthStatus()),
		SyncStatus:   FromArgoCDAppSyncStatusProto(a.GetSyncStatus()),
	}
}

func FromArgoCDAppHealthStatusProto(
	a *v1alpha1.ArgoCDAppHealthStatus,
) kargoapi.ArgoCDAppHealthStatus {
	return kargoapi.ArgoCDAppHealthStatus{
		Status:  kargoapi.ArgoCDAppHealthState(a.GetStatus()),
		Message: a.GetMessage(),
	}
}

func FromArgoCDAppSyncStatusProto(
	a *v1alpha1.ArgoCDAppSyncStatus,
) kargoapi.ArgoCDAppSyncStatus {
	return kargoapi.ArgoCDAppSyncStatus{
		Status:    kargoapi.ArgoCDAppSyncState(a.GetStatus()),
		Revision:  a.GetRevision(),
		Revisions: a.GetRevisions(),
	}
}

func FromSubscriptionsProto(s *v1alpha1.Subscriptions) *kargoapi.Subscriptions {
	if s == nil {
		return nil
	}
	upstreamStages := make([]kargoapi.StageSubscription, len(s.GetUpstreamStages()))
	for idx, stage := range s.GetUpstreamStages() {
		upstreamStages[idx] = *FromStageSubscriptionProto(stage)
	}
	return &kargoapi.Subscriptions{
		Warehouse:      s.GetWarehouse(),
		UpstreamStages: upstreamStages,
	}
}

func FromRepoSubscriptionProto(s *v1alpha1.RepoSubscription) *kargoapi.RepoSubscription {
	if s == nil {
		return nil
	}
	return &kargoapi.RepoSubscription{
		Git:   FromGitSubscriptionProto(s.Git),
		Image: FromImageSubscriptionProto(s.Image),
		Chart: FromChartSubscriptionProto(s.Chart),
	}
}

func FromGitSubscriptionProto(s *v1alpha1.GitSubscription) *kargoapi.GitSubscription {
	if s == nil {
		return nil
	}
	return &kargoapi.GitSubscription{
		RepoURL: s.GetRepoUrl(),
		Branch:  s.GetBranch(),
	}
}

func FromImageSubscriptionProto(s *v1alpha1.ImageSubscription) *kargoapi.ImageSubscription {
	if s == nil {
		return nil
	}
	return &kargoapi.ImageSubscription{
		RepoURL:          s.GetRepoUrl(),
		UpdateStrategy:   kargoapi.ImageUpdateStrategy(s.GetUpdateStrategy()),
		SemverConstraint: s.GetSemverConstraint(),
		AllowTags:        s.GetAllowTags(),
		IgnoreTags:       s.GetIgnoreTags(),
		Platform:         s.GetPlatform(),
	}
}

func FromChartSubscriptionProto(s *v1alpha1.ChartSubscription) *kargoapi.ChartSubscription {
	if s == nil {
		return nil
	}
	return &kargoapi.ChartSubscription{
		RegistryURL:      s.GetRegistryUrl(),
		Name:             s.GetName(),
		SemverConstraint: s.GetSemverConstraint(),
	}
}

func FromPromotionMechanismsProto(m *v1alpha1.PromotionMechanisms) *kargoapi.PromotionMechanisms {
	if m == nil {
		return nil
	}
	gitUpdates := make([]kargoapi.GitRepoUpdate, len(m.GetGitRepoUpdates()))
	for idx, git := range m.GetGitRepoUpdates() {
		gitUpdates[idx] = *FromGitRepoUpdateProto(git)
	}
	argoUpdates := make([]kargoapi.ArgoCDAppUpdate, len(m.GetArgocdAppUpdates()))
	for idx, argo := range m.GetArgocdAppUpdates() {
		argoUpdates[idx] = *FromArgoCDAppUpdatesProto(argo)
	}
	return &kargoapi.PromotionMechanisms{
		GitRepoUpdates:   gitUpdates,
		ArgoCDAppUpdates: argoUpdates,
	}
}

func FromGitRepoUpdateProto(u *v1alpha1.GitRepoUpdate) *kargoapi.GitRepoUpdate {
	if u == nil {
		return nil
	}
	return &kargoapi.GitRepoUpdate{
		RepoURL:     u.GetRepoUrl(),
		ReadBranch:  u.GetReadBranch(),
		WriteBranch: u.GetWriteBranch(),
		Render:      FromKargoRenderPromotionMechanismProto(u.GetRender()),
		Kustomize:   FromKustomizePromotionMechanismProto(u.GetKustomize()),
		Helm:        FromHelmPromotionMechanismProto(u.GetHelm()),
	}
}

func FromKargoRenderPromotionMechanismProto(
	m *v1alpha1.KargoRenderPromotionMechanism,
) *kargoapi.KargoRenderPromotionMechanism {
	if m == nil {
		return nil
	}
	return &kargoapi.KargoRenderPromotionMechanism{}
}

func FromKustomizePromotionMechanismProto(
	m *v1alpha1.KustomizePromotionMechanism,
) *kargoapi.KustomizePromotionMechanism {
	if m == nil {
		return nil
	}
	images := make([]kargoapi.KustomizeImageUpdate, len(m.GetImages()))
	for idx, image := range m.GetImages() {
		images[idx] = *FromKustomizeImageUpdateProto(image)
	}
	return &kargoapi.KustomizePromotionMechanism{
		Images: images,
	}
}

func FromKustomizeImageUpdateProto(u *v1alpha1.KustomizeImageUpdate) *kargoapi.KustomizeImageUpdate {
	if u == nil {
		return nil
	}
	return &kargoapi.KustomizeImageUpdate{
		Image: u.GetImage(),
		Path:  u.GetPath(),
	}
}

func FromHelmPromotionMechanismProto(
	m *v1alpha1.HelmPromotionMechanism,
) *kargoapi.HelmPromotionMechanism {
	if m == nil {
		return nil
	}
	images := make([]kargoapi.HelmImageUpdate, len(m.GetImages()))
	for idx, image := range m.GetImages() {
		images[idx] = *FromHelmImageUpdateProto(image)
	}
	charts := make([]kargoapi.HelmChartDependencyUpdate, len(m.GetCharts()))
	for idx, chart := range m.GetCharts() {
		charts[idx] = *FromHelmChartDependencyUpdateProto(chart)
	}
	return &kargoapi.HelmPromotionMechanism{
		Images: images,
		Charts: charts,
	}
}

func FromHelmImageUpdateProto(u *v1alpha1.HelmImageUpdate) *kargoapi.HelmImageUpdate {
	if u == nil {
		return nil
	}
	return &kargoapi.HelmImageUpdate{
		Image:          u.GetImage(),
		ValuesFilePath: u.GetValuesFilePath(),
		Key:            u.GetKey(),
		Value:          kargoapi.ImageUpdateValueType(u.GetValue()),
	}
}

func FromHelmChartDependencyUpdateProto(
	u *v1alpha1.HelmChartDependencyUpdate,
) *kargoapi.HelmChartDependencyUpdate {
	if u == nil {
		return nil
	}
	return &kargoapi.HelmChartDependencyUpdate{
		RegistryURL: u.GetRegistryUrl(),
		Name:        u.GetName(),
		ChartPath:   u.GetChartPath(),
	}
}

func FromArgoCDAppUpdatesProto(u *v1alpha1.ArgoCDAppUpdate) *kargoapi.ArgoCDAppUpdate {
	if u == nil {
		return nil
	}
	sourceUpdates := make([]kargoapi.ArgoCDSourceUpdate, len(u.GetSourceUpdates()))
	for idx, update := range u.GetSourceUpdates() {
		sourceUpdates[idx] = *FromArgoCDSourceUpdateProto(update)
	}
	return &kargoapi.ArgoCDAppUpdate{
		AppName:       u.GetAppName(),
		AppNamespace:  u.GetAppNamespace(),
		SourceUpdates: sourceUpdates,
	}
}

func FromArgoCDSourceUpdateProto(u *v1alpha1.ArgoCDSourceUpdate) *kargoapi.ArgoCDSourceUpdate {
	if u == nil {
		return nil
	}
	return &kargoapi.ArgoCDSourceUpdate{
		RepoURL:              u.GetRepoUrl(),
		Chart:                u.GetChart(),
		UpdateTargetRevision: u.GetUpdateTargetRevision(),
		Kustomize:            FromArgoCDKustomizeProto(u.GetKustomize()),
		Helm:                 FromArgoCDHelm(u.GetHelm()),
	}
}

func FromArgoCDKustomizeProto(k *v1alpha1.ArgoCDKustomize) *kargoapi.ArgoCDKustomize {
	if k == nil {
		return nil
	}
	return &kargoapi.ArgoCDKustomize{
		Images: k.GetImages(),
	}
}

func FromArgoCDHelm(h *v1alpha1.ArgoCDHelm) *kargoapi.ArgoCDHelm {
	if h == nil {
		return nil
	}
	images := make([]kargoapi.ArgoCDHelmImageUpdate, len(h.GetImages()))
	for idx, image := range h.GetImages() {
		images[idx] = *FromArgoCDHelmImageUpdateProto(image)
	}
	return &kargoapi.ArgoCDHelm{
		Images: images,
	}
}

func FromArgoCDHelmImageUpdateProto(u *v1alpha1.ArgoCDHelmImageUpdate) *kargoapi.ArgoCDHelmImageUpdate {
	if u == nil {
		return nil
	}
	return &kargoapi.ArgoCDHelmImageUpdate{
		Image: u.GetImage(),
		Key:   u.GetKey(),
		Value: kargoapi.ImageUpdateValueType(u.GetValue()),
	}
}

func FromStageSubscriptionProto(s *v1alpha1.StageSubscription) *kargoapi.StageSubscription {
	if s == nil {
		return nil
	}
	return &kargoapi.StageSubscription{
		Name: s.GetName(),
	}
}

func FromPromotionProto(p *v1alpha1.Promotion) *kargoapi.Promotion {
	if p == nil {
		return nil
	}
	var status kargoapi.PromotionStatus
	if p.GetStatus() != nil {
		status = *FromPromotionStatusProto(p.GetStatus())
	}
	var objectMeta kubemetav1.ObjectMeta
	if p.GetMetadata() != nil {
		objectMeta = *typesmetav1.FromObjectMetaProto(p.GetMetadata())
	}
	return &kargoapi.Promotion{
		TypeMeta: kubemetav1.TypeMeta{
			APIVersion: kargoapi.GroupVersion.String(),
			Kind:       "Promotion",
		},
		ObjectMeta: objectMeta,
		Spec:       FromPromotionSpecProto(p.GetSpec()),
		Status:     status,
	}
}

func FromPromotionSpecProto(s *v1alpha1.PromotionSpec) *kargoapi.PromotionSpec {
	if s == nil {
		return nil
	}
	return &kargoapi.PromotionSpec{
		Stage:   s.GetStage(),
		Freight: s.GetFreight(),
	}
}

func FromPromotionStatusProto(s *v1alpha1.PromotionStatus) *kargoapi.PromotionStatus {
	if s == nil {
		return nil
	}
	return &kargoapi.PromotionStatus{
		Phase: kargoapi.PromotionPhase(s.GetPhase()),
		Error: s.GetError(),
	}
}

func FromPromotionPolicyProto(p *v1alpha1.PromotionPolicy) *kargoapi.PromotionPolicy {
	if p == nil {
		return nil
	}
	var objectMeta kubemetav1.ObjectMeta
	if p.GetMetadata() != nil {
		objectMeta = *typesmetav1.FromObjectMetaProto(p.GetMetadata())
	}
	return &kargoapi.PromotionPolicy{
		TypeMeta: kubemetav1.TypeMeta{
			APIVersion: kargoapi.GroupVersion.String(),
			Kind:       "PromotionPolicy",
		},
		ObjectMeta:          objectMeta,
		Stage:               p.GetStage(),
		EnableAutoPromotion: p.GetEnableAutoPromotion(),
	}
}

func ToStageProto(e kargoapi.Stage) *v1alpha1.Stage {
	// Status
	var currentFreight *v1alpha1.SimpleFreight
	if e.Status.CurrentFreight != nil {
		currentFreight = ToSimpleFreightProto(*e.Status.CurrentFreight, nil)
	}
	history := make([]*v1alpha1.SimpleFreight, len(e.Status.History))
	for idx := range e.Status.History {
		history[idx] = ToSimpleFreightProto(e.Status.History[idx], nil)
	}
	var health *v1alpha1.Health
	if e.Status.Health != nil {
		health = ToHealthProto(*e.Status.Health)
	}

	metadata := e.ObjectMeta.DeepCopy()
	metadata.SetManagedFields(nil)

	var promotionMechanisms *v1alpha1.PromotionMechanisms
	if e.Spec.PromotionMechanisms != nil {
		promotionMechanisms = ToPromotionMechanismsProto(*e.Spec.PromotionMechanisms)
	}
	var currentPromotion *v1alpha1.PromotionInfo
	if e.Status.CurrentPromotion != nil {
		sf := kargoapi.SimpleFreight{
			ID:      e.Status.CurrentPromotion.Freight.ID,
			Commits: e.Status.CurrentPromotion.Freight.Commits,
			Images:  e.Status.CurrentPromotion.Freight.Images,
			Charts:  e.Status.CurrentPromotion.Freight.Charts,
		}
		currentPromotion = &v1alpha1.PromotionInfo{
			Name:    e.Status.CurrentPromotion.Name,
			Freight: ToSimpleFreightProto(sf, nil),
		}
	}
	return &v1alpha1.Stage{
		ApiVersion: e.APIVersion,
		Kind:       e.Kind,
		Metadata:   typesmetav1.ToObjectMetaProto(*metadata),
		Spec: &v1alpha1.StageSpec{
			Subscriptions:       ToSubscriptionsProto(*e.Spec.Subscriptions),
			PromotionMechanisms: promotionMechanisms,
		},
		Status: &v1alpha1.StageStatus{
			CurrentFreight:   currentFreight,
			CurrentPromotion: currentPromotion,
			History:          history,
			Health:           health,
			Error:            e.Status.Error,
		},
	}
}

func ToRepoSubscriptionProto(s kargoapi.RepoSubscription) *v1alpha1.RepoSubscription {
	var git *v1alpha1.GitSubscription
	if s.Git != nil {
		git = ToGitSubscriptionProto(*s.Git)
	}
	var image *v1alpha1.ImageSubscription
	if s.Image != nil {
		image = ToImageSubscriptionProto(*s.Image)
	}
	var chart *v1alpha1.ChartSubscription
	if s.Chart != nil {
		chart = ToChartSubscriptionProto(*s.Chart)
	}
	return &v1alpha1.RepoSubscription{
		Git:   git,
		Image: image,
		Chart: chart,
	}
}

func ToSubscriptionsProto(s kargoapi.Subscriptions) *v1alpha1.Subscriptions {
	upstreamStages := make([]*v1alpha1.StageSubscription, len(s.UpstreamStages))
	for idx := range s.UpstreamStages {
		upstreamStages[idx] = ToStageSubscriptionProto(s.UpstreamStages[idx])
	}
	return &v1alpha1.Subscriptions{
		Warehouse:      s.Warehouse,
		UpstreamStages: upstreamStages,
	}
}

func ToGitSubscriptionProto(g kargoapi.GitSubscription) *v1alpha1.GitSubscription {
	return &v1alpha1.GitSubscription{
		RepoUrl: g.RepoURL,
		Branch:  g.Branch,
	}
}

func ToImageSubscriptionProto(i kargoapi.ImageSubscription) *v1alpha1.ImageSubscription {
	return &v1alpha1.ImageSubscription{
		RepoUrl:          i.RepoURL,
		UpdateStrategy:   string(i.UpdateStrategy),
		SemverConstraint: proto.String(i.SemverConstraint),
		AllowTags:        proto.String(i.AllowTags),
		IgnoreTags:       i.IgnoreTags,
		Platform:         proto.String(i.Platform),
	}
}

func ToChartSubscriptionProto(c kargoapi.ChartSubscription) *v1alpha1.ChartSubscription {
	return &v1alpha1.ChartSubscription{
		RegistryUrl:      c.RegistryURL,
		Name:             proto.String(c.Name),
		SemverConstraint: proto.String(c.SemverConstraint),
	}
}

func ToStageSubscriptionProto(e kargoapi.StageSubscription) *v1alpha1.StageSubscription {
	return &v1alpha1.StageSubscription{
		Name: e.Name,
	}
}

func ToPromotionMechanismsProto(p kargoapi.PromotionMechanisms) *v1alpha1.PromotionMechanisms {
	gitRepoUpdates := make([]*v1alpha1.GitRepoUpdate, len(p.GitRepoUpdates))
	for idx := range p.GitRepoUpdates {
		gitRepoUpdates[idx] = ToGitRepoUpdateProto(p.GitRepoUpdates[idx])
	}
	argoCDAppUpdates := make([]*v1alpha1.ArgoCDAppUpdate, len(p.ArgoCDAppUpdates))
	for idx := range p.ArgoCDAppUpdates {
		argoCDAppUpdates[idx] = ToArgoCDAppUpdateProto(p.ArgoCDAppUpdates[idx])
	}
	return &v1alpha1.PromotionMechanisms{
		GitRepoUpdates:   gitRepoUpdates,
		ArgocdAppUpdates: argoCDAppUpdates,
	}
}

func ToGitRepoUpdateProto(g kargoapi.GitRepoUpdate) *v1alpha1.GitRepoUpdate {
	var render *v1alpha1.KargoRenderPromotionMechanism
	if g.Render != nil {
		render = ToKargoRenderPromotionMechanismProto(*g.Render)
	}
	var kustomize *v1alpha1.KustomizePromotionMechanism
	if g.Kustomize != nil {
		kustomize = ToKustomizePromotionMechanismProto(*g.Kustomize)
	}
	var helm *v1alpha1.HelmPromotionMechanism
	if g.Helm != nil {
		helm = ToHelmPromotionMechanismProto(*g.Helm)
	}
	return &v1alpha1.GitRepoUpdate{
		RepoUrl:     g.RepoURL,
		ReadBranch:  proto.String(g.ReadBranch),
		WriteBranch: g.WriteBranch,
		Render:      render,
		Kustomize:   kustomize,
		Helm:        helm,
	}
}

func ToKargoRenderPromotionMechanismProto(
	_ kargoapi.KargoRenderPromotionMechanism,
) *v1alpha1.KargoRenderPromotionMechanism {
	return &v1alpha1.KargoRenderPromotionMechanism{}
}

func ToKustomizePromotionMechanismProto(
	k kargoapi.KustomizePromotionMechanism,
) *v1alpha1.KustomizePromotionMechanism {
	images := make([]*v1alpha1.KustomizeImageUpdate, len(k.Images))
	for idx := range k.Images {
		images[idx] = ToKustomizeImageUpdateProto(k.Images[idx])
	}
	return &v1alpha1.KustomizePromotionMechanism{
		Images: images,
	}
}

func ToKustomizeImageUpdateProto(k kargoapi.KustomizeImageUpdate) *v1alpha1.KustomizeImageUpdate {
	return &v1alpha1.KustomizeImageUpdate{
		Image: k.Image,
		Path:  k.Path,
	}
}

func ToHelmPromotionMechanismProto(h kargoapi.HelmPromotionMechanism) *v1alpha1.HelmPromotionMechanism {
	images := make([]*v1alpha1.HelmImageUpdate, len(h.Images))
	for idx := range h.Images {
		images[idx] = ToHelmImageUpdateProto(h.Images[idx])
	}
	charts := make([]*v1alpha1.HelmChartDependencyUpdate, len(h.Charts))
	for idx := range h.Charts {
		charts[idx] = ToHelmChartDependencyUpdateProto(h.Charts[idx])
	}
	return &v1alpha1.HelmPromotionMechanism{
		Images: images,
		Charts: charts,
	}
}

func ToHelmImageUpdateProto(h kargoapi.HelmImageUpdate) *v1alpha1.HelmImageUpdate {
	return &v1alpha1.HelmImageUpdate{
		Image:          h.Image,
		ValuesFilePath: h.ValuesFilePath,
		Key:            h.Key,
		Value:          string(h.Value),
	}
}

func ToHelmChartDependencyUpdateProto(h kargoapi.HelmChartDependencyUpdate) *v1alpha1.HelmChartDependencyUpdate {
	return &v1alpha1.HelmChartDependencyUpdate{
		RegistryUrl: h.RegistryURL,
		Name:        h.Name,
		ChartPath:   h.ChartPath,
	}
}

func ToArgoCDAppUpdateProto(h kargoapi.ArgoCDAppUpdate) *v1alpha1.ArgoCDAppUpdate {
	sourceUpdates := make([]*v1alpha1.ArgoCDSourceUpdate, len(h.SourceUpdates))
	for idx := range h.SourceUpdates {
		sourceUpdates[idx] = ToArgoCDSourceUpdateProto(h.SourceUpdates[idx])
	}
	return &v1alpha1.ArgoCDAppUpdate{
		AppName:       h.AppName,
		AppNamespace:  proto.String(h.AppNamespace),
		SourceUpdates: sourceUpdates,
	}
}

func ToArgoCDSourceUpdateProto(a kargoapi.ArgoCDSourceUpdate) *v1alpha1.ArgoCDSourceUpdate {
	var kustomize *v1alpha1.ArgoCDKustomize
	if a.Kustomize != nil {
		kustomize = ToArgoCDKustomizeProto(*a.Kustomize)
	}
	var helm *v1alpha1.ArgoCDHelm
	if a.Helm != nil {
		helm = ToArgoCDHelmProto(*a.Helm)
	}
	return &v1alpha1.ArgoCDSourceUpdate{
		RepoUrl:              a.RepoURL,
		Chart:                proto.String(a.Chart),
		UpdateTargetRevision: proto.Bool(a.UpdateTargetRevision),
		Kustomize:            kustomize,
		Helm:                 helm,
	}
}

func ToArgoCDKustomizeProto(a kargoapi.ArgoCDKustomize) *v1alpha1.ArgoCDKustomize {
	return &v1alpha1.ArgoCDKustomize{
		Images: a.Images,
	}
}

func ToArgoCDHelmProto(a kargoapi.ArgoCDHelm) *v1alpha1.ArgoCDHelm {
	images := make([]*v1alpha1.ArgoCDHelmImageUpdate, len(a.Images))
	for idx := range images {
		images[idx] = ToArgoCDHelmImageUpdateProto(a.Images[idx])
	}
	return &v1alpha1.ArgoCDHelm{
		Images: images,
	}
}

func ToArgoCDHelmImageUpdateProto(a kargoapi.ArgoCDHelmImageUpdate) *v1alpha1.ArgoCDHelmImageUpdate {
	return &v1alpha1.ArgoCDHelmImageUpdate{
		Image: a.Image,
		Key:   a.Key,
		Value: string(a.Value),
	}
}

func ToFreightProto(f kargoapi.Freight) *v1alpha1.Freight {
	metadata := f.ObjectMeta.DeepCopy()
	metadata.SetManagedFields(nil)
	commits := make([]*v1alpha1.GitCommit, len(f.Commits))
	for idx := range f.Commits {
		commits[idx] = ToGitCommitProto(f.Commits[idx])
	}
	images := make([]*v1alpha1.Image, len(f.Images))
	for idx := range f.Images {
		images[idx] = ToImageProto(f.Images[idx])
	}
	charts := make([]*v1alpha1.Chart, len(f.Charts))
	for idx := range f.Charts {
		charts[idx] = ToChartProto(f.Charts[idx])
	}
	qualifications :=
		make(map[string]*v1alpha1.Qualification, len(f.Status.Qualifications))
	for stageName := range f.Status.Qualifications {
		qualifications[stageName] = &v1alpha1.Qualification{}
	}
	return &v1alpha1.Freight{
		ApiVersion: f.APIVersion,
		Kind:       f.Kind,
		Id:         f.ID,
		Images:     images,
		Charts:     charts,
		Commits:    commits,
		Metadata:   typesmetav1.ToObjectMetaProto(*metadata),
		Status: &v1alpha1.FreightStatus{
			Qualifications: qualifications,
		},
	}
}

func ToSimpleFreightProto(s kargoapi.SimpleFreight, firstSeen *time.Time) *v1alpha1.SimpleFreight {
	var firstSeenProto *timestamppb.Timestamp
	if firstSeen != nil {
		firstSeenProto = timestamppb.New(*firstSeen)
	}
	commits := make([]*v1alpha1.GitCommit, len(s.Commits))
	for idx := range s.Commits {
		commits[idx] = ToGitCommitProto(s.Commits[idx])
	}
	images := make([]*v1alpha1.Image, len(s.Images))
	for idx := range s.Images {
		images[idx] = ToImageProto(s.Images[idx])
	}
	charts := make([]*v1alpha1.Chart, len(s.Charts))
	for idx := range s.Charts {
		charts[idx] = ToChartProto(s.Charts[idx])
	}
	return &v1alpha1.SimpleFreight{
		Id:        s.ID,
		FirstSeen: firstSeenProto,
		Commits:   commits,
		Images:    images,
		Charts:    charts,
	}
}

func ToWarehouseProto(w kargoapi.Warehouse) *v1alpha1.Warehouse {
	subscriptions := make([]*v1alpha1.RepoSubscription, len(w.Spec.Subscriptions))
	for idx, subscription := range w.Spec.Subscriptions {
		subscriptions[idx] = ToRepoSubscriptionProto(subscription)
	}
	var status *v1alpha1.WarehouseStatus
	if w.GetStatus() != nil {
		status = &v1alpha1.WarehouseStatus{
			Error:              w.GetStatus().Error,
			ObservedGeneration: w.GetStatus().ObservedGeneration,
		}
	}
	return &v1alpha1.Warehouse{
		ApiVersion: w.APIVersion,
		Kind:       w.Kind,
		Metadata:   typesmetav1.ToObjectMetaProto(w.ObjectMeta),
		Spec: &v1alpha1.WarehouseSpec{
			Subscriptions: subscriptions,
		},
		Status: status,
	}
}

func ToGitCommitProto(g kargoapi.GitCommit) *v1alpha1.GitCommit {
	return &v1alpha1.GitCommit{
		RepoUrl:           g.RepoURL,
		Id:                g.ID,
		Branch:            g.Branch,
		HealthCheckCommit: proto.String(g.HealthCheckCommit),
		Message:           g.Message,
		Author:            g.Author,
	}
}

func ToImageProto(i kargoapi.Image) *v1alpha1.Image {
	return &v1alpha1.Image{
		RepoUrl: i.RepoURL,
		Tag:     i.Tag,
	}
}

func ToChartProto(c kargoapi.Chart) *v1alpha1.Chart {
	return &v1alpha1.Chart{
		RegistryUrl: c.RegistryURL,
		Name:        c.Name,
		Version:     c.Version,
	}
}

func ToHealthProto(h kargoapi.Health) *v1alpha1.Health {
	argocdAppStates := make([]*v1alpha1.ArgoCDAppState, len(h.ArgoCDApps))
	for i, argocdAppState := range h.ArgoCDApps {
		argocdAppStates[i] = ToArgoCDAppStateProto(argocdAppState)
	}
	return &v1alpha1.Health{
		Status:     string(h.Status),
		Issues:     h.Issues,
		ArgocdApps: argocdAppStates,
	}
}

func ToArgoCDAppStateProto(
	a kargoapi.ArgoCDAppStatus,
) *v1alpha1.ArgoCDAppState {
	return &v1alpha1.ArgoCDAppState{
		Name:         a.Name,
		Namespace:    a.Namespace,
		HealthStatus: ToArgoCDAppHealthStatusProto(a.HealthStatus),
		SyncStatus:   ToArgoCDAppSyncStatusProto(a.SyncStatus),
	}
}

func ToArgoCDAppHealthStatusProto(
	a kargoapi.ArgoCDAppHealthStatus,
) *v1alpha1.ArgoCDAppHealthStatus {
	return &v1alpha1.ArgoCDAppHealthStatus{
		Status:  string(a.Status),
		Message: a.Message,
	}
}

func ToArgoCDAppSyncStatusProto(
	a kargoapi.ArgoCDAppSyncStatus,
) *v1alpha1.ArgoCDAppSyncStatus {
	return &v1alpha1.ArgoCDAppSyncStatus{
		Status:    string(a.Status),
		Revision:  a.Revision,
		Revisions: a.Revisions,
	}
}

func ToPromotionProto(p kargoapi.Promotion) *v1alpha1.Promotion {
	metadata := p.ObjectMeta.DeepCopy()
	metadata.SetManagedFields(nil)

	return &v1alpha1.Promotion{
		ApiVersion: p.APIVersion,
		Kind:       p.Kind,
		Metadata:   typesmetav1.ToObjectMetaProto(*metadata),
		Spec: &v1alpha1.PromotionSpec{
			Stage:   p.Spec.Stage,
			Freight: p.Spec.Freight,
		},
		Status: &v1alpha1.PromotionStatus{
			Phase: string(p.Status.Phase),
			Error: p.Status.Error,
		},
	}
}

func ToPromotionPolicyProto(p kargoapi.PromotionPolicy) *v1alpha1.PromotionPolicy {
	metadata := p.ObjectMeta.DeepCopy()
	metadata.SetManagedFields(nil)

	return &v1alpha1.PromotionPolicy{
		ApiVersion:          p.APIVersion,
		Kind:                p.Kind,
		Metadata:            typesmetav1.ToObjectMetaProto(*metadata),
		Stage:               p.Stage,
		EnableAutoPromotion: p.EnableAutoPromotion,
	}
}

func ToVersionProto(v version.Version) *svcv1alpha1.VersionInfo {
	return &svcv1alpha1.VersionInfo{
		Version:      v.Version,
		GitCommit:    v.GitCommit,
		GitTreeDirty: v.GitTreeDirty,
		BuildTime:    timestamppb.New(v.BuildDate),
		GoVersion:    v.GoVersion,
		Compiler:     v.Compiler,
		Platform:     v.Platform,
	}
}
